package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/csothen/tmdei-project/internal/db"
	"github.com/csothen/tmdei-project/internal/models"
	"github.com/csothen/tmdei-project/internal/models/dtos"
	"github.com/csothen/tmdei-project/internal/terraform"
	"github.com/csothen/tmdei-project/internal/utils"
)

// ReadAll retrieves all deployments
func (s *Service) ReadAllDeployments(ctx context.Context) []*models.Deployment {
	dbds, err := s.repo.GetAllDeployments(ctx)
	if err != nil {
		return nil
	}
	deployments := make([]*models.Deployment, len(dbds))
	for i, dbd := range dbds {
		var deployment *models.Deployment
		deployment.FromDB(dbd)
		deployments[i] = deployment
	}
	return deployments
}

// ReadOne retrieves a deployment that matches the canonical
func (s *Service) ReadOneDeployment(ctx context.Context, dcan string) (*models.Deployment, error) {
	dbd, err := s.repo.GetDeploymentByCanonical(ctx, dcan)
	if err != nil {
		return nil, fmt.Errorf("deployment %s not found: %w", dcan, err)
	}

	var deployment *models.Deployment
	deployment.FromDB(dbd)
	return deployment, nil
}

// Create creates a deployment with the inputs given
func (s *Service) CreateDeployment(ctx context.Context, nds *dtos.NewDeployments) (deployments []*dtos.CreatedDeployment, warnings, errors []error) {
	warnings = make([]error, 0)
	errors = make([]error, 0)

	// Validate the requested deployments
	for _, nd := range nds.Deployments {
		uc := nd.UseCase
		if _, err := s.repo.GetUseCaseConfiguration(ctx, uc); err != nil {
			errors = append(errors, fmt.Errorf("could not find usecase %s", uc))
			continue
		}

		for _, ns := range nd.Services {
			st, err := models.TypeString(ns.Service)
			if err != nil {
				errors = append(errors, fmt.Errorf("%s is not a valid service", ns.Service))
				continue
			}
			if _, err := s.repo.GetServiceConfiguration(ctx, uc, uint(st)); err != nil {
				errors = append(errors, fmt.Errorf("could not find configuration for service %s on usecase %s", ns.Service, uc))
				continue
			}
		}
	}

	// If there were errors then we return the
	// errors and don't process the deployments
	if len(errors) > 0 {
		return nil, nil, errors
	}

	// create the Terraform worker which will do the deployments
	tfw := terraform.NewWorker(s.config.TerraformExecPath)

	// find the template files that will be used to generate the
	// actual deployment files
	templatesDir, err := utils.BuildTemplatesFolderPath()
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build templates directory path: %w", err))
	}

	dirs, err := os.ReadDir(templatesDir)
	if err != nil {
		errors = append(errors, fmt.Errorf("could not read all template directories: %w", err))
		return nil, nil, errors
	}

	// load the template files paths in a map
	templates := make(map[string][]string)
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		serviceTemplateDir := path.Join(templatesDir, dir.Name())
		fis, err := ioutil.ReadDir(serviceTemplateDir)
		if err != nil {
			errors = append(errors, fmt.Errorf("could not read all files in dir %s: %w", serviceTemplateDir, err))
			continue
		}

		templates[dir.Name()] = make([]string, 0)
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			filePaths := templates[dir.Name()]
			filePaths = append(filePaths, path.Join(serviceTemplateDir, fi.Name()))
			templates[dir.Name()] = filePaths
		}
	}

	// if for some reson we failed to find the files
	// we return the errors
	if len(errors) > 0 {
		return nil, nil, errors
	}

	// For each individual deployment persist its information
	// and start a goroutine that will do the deployment logic
	wg := sync.WaitGroup{}
	deployments = make([]*dtos.CreatedDeployment, 0)
	for _, nd := range nds.Deployments {
		for _, ns := range nd.Services {
			for i := 0; i < ns.Count; i++ {
				canonical := utils.BuildDeploymentCanonical(nd.UseCase, ns.Service, i+1)
				// was previously validated
				st, _ := models.TypeString(ns.Service)

				d, warning := s.persistDeployment(ctx, canonical, st, nds.CallbackURL)
				if warning != nil {
					warnings = append(warnings, warning)
					continue
				}

				cd := &dtos.CreatedDeployment{
					Canonical: canonical,
					State:     d.State.String(),
					Type:      d.Type.String(),
				}
				deployments = append(deployments, cd)

				intpl := utils.Interpolator{
					Name:  canonical,
					Count: 1,
				}

				// we start goroutines that will do the actual deployments
				// the deployment consists of:
				// 1. Read the template files
				// 2. Replace what needs to be replaced in the template file
				// 3. Persist the resulting file
				// 4. Run the terraform using the terraform worker with the resulting files
				wg.Add(1)
				go func() {
					defer wg.Done()

					pathToDir, err := utils.BuildDeploymentFolderPath(cd.Canonical)
					if err != nil {
						errors = append(errors, fmt.Errorf("could not build deployments path: %w", err))
						return
					}

					err = os.MkdirAll(pathToDir, 0755)
					if err != nil {
						errors = append(errors, fmt.Errorf("error creating directory for deployment %s: %w", cd.Canonical, err))
						return
					}

					filePaths, ok := templates[ns.Service]
					if !ok {
						errors = append(errors, fmt.Errorf("no template founds for the service %s", ns.Service))
						return
					}

					for _, fp := range filePaths {
						f, err := os.Open(fp)
						if err != nil {
							errors = append(errors, fmt.Errorf("could not open file %s: %w", fp, err))
							return
						}
						defer f.Close()

						fcontents, err := ioutil.ReadAll(f)
						if err != nil {
							errors = append(errors, fmt.Errorf("could not read file %s's contents: %w", fp, err))
							return
						}

						fcontents = intpl.Inteprolate(fcontents)
						filePath := path.Join(pathToDir, path.Base(fp))
						cf, err := os.Create(filePath)
						if err != nil {
							errors = append(errors, fmt.Errorf("could not create file %s: %w", filePath, err))
							return
						}
						defer cf.Close()

						_, err = cf.Write(fcontents)
						if err != nil {
							errors = append(errors, fmt.Errorf("could not write to file %s: %w", filePath, err))
						}
					}

					err = tfw.Deploy(pathToDir)
					if err != nil {
						errors = append(errors, fmt.Errorf("error executing deployment %s: %w", cd.Canonical, err))
						return
					}

					outputs, err := tfw.Outputs(pathToDir)
					if err != nil {
						errors = append(errors, fmt.Errorf("error retrieving the public IP of the deployment %s: %w", cd.Canonical, err))
						return
					}

					deploymentURL := outputs["public_ips"]
					dbd := d.ToDB()
					dbd.URL = deploymentURL
					err = s.repo.UpdateDeployment(ctx, *dbd)
					if err != nil {
						errors = append(errors, fmt.Errorf("error updating deployment %s's URL: %w", cd.Canonical, err))
						return
					}
				}()
			}
		}
	}

	// Wait for all the deployment tasks to be over
	// and check for errors in order for them to be logged
	go func() {
		wg.Wait()
		for _, err := range errors {
			log.Println(err)
		}
	}()
	return deployments, warnings, nil
}

func (s *Service) persistDeployment(ctx context.Context, canonical string, st models.Type, callbackURL string) (*models.Deployment, error) {
	dbd, err := s.repo.CreateDeployment(ctx, db.Deployment{
		Canonical:   canonical,
		State:       uint(models.Pending),
		Type:        uint(st),
		CallbackURL: callbackURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to persist the deployment %s: %w", canonical, err)
	}
	var deployment *models.Deployment
	deployment.FromDB(dbd)

	return deployment, nil
}

// Update updates a deployment with the given canonical
func (s *Service) UpdateDeployment(ctx context.Context, dcan string, d *models.Deployment) error {
	err := s.repo.UpdateDeployment(ctx, *d.ToDB())
	if err != nil {
		return fmt.Errorf("failed to update the deployment %s: %w", dcan, err)
	}
	return nil
}

// Delete deletes a deployment with the given canonical
func (s *Service) DeleteDeployment(ctx context.Context, dcan string) error {
	err := s.repo.DeleteDeployment(ctx, dcan)
	if err != nil {
		return fmt.Errorf("failed to delete deployment %s: %w", dcan, err)
	}
	return nil
}
