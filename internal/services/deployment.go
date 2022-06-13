package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/csothen/lift/internal/db"
	"github.com/csothen/lift/internal/fetcher"
	"github.com/csothen/lift/internal/fetcher/sonarqube"
	"github.com/csothen/lift/internal/models"
	"github.com/csothen/lift/internal/models/dtos"
	"github.com/csothen/lift/internal/terraform"
	"github.com/csothen/lift/internal/utils"
)

// ReadAll retrieves all deployments
func (s *Service) ReadAllDeployments(ctx context.Context) []*models.Deployment {
	dbds, err := s.repo.GetAllDeployments(ctx)
	if err != nil {
		return nil
	}
	deployments := make([]*models.Deployment, len(dbds))
	for i, dbd := range dbds {
		deployment := &models.Deployment{}
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

	deployment := &models.Deployment{}
	deployment.FromDB(dbd)
	return deployment, nil
}

// Create creates a deployment with the inputs given
func (s *Service) CreateDeployment(ctx context.Context, nds *dtos.NewDeployments) (deployments []*dtos.CreatedDeployment, warnings, errors []error) {
	warnings = make([]error, 0)
	errors = make([]error, 0)

	configurations := make(map[string]*models.ServiceConfiguration)
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
			dbsc, err := s.repo.GetServiceConfiguration(ctx, uc, uint(st))
			if err != nil {
				errors = append(errors, fmt.Errorf("could not find configuration for service %s on usecase %s", ns.Service, uc))
				continue
			}

			// load the configuration for the specific use case and service type
			sc := &models.ServiceConfiguration{}
			sc.FromDB(dbsc)
			key := fmt.Sprintf("%s:%s", uc, ns.Service)
			configurations[key] = sc
		}
	}

	// If there were errors then we return the
	// errors and don't process the deployments
	if len(errors) > 0 {
		return nil, nil, errors
	}

	// Load the fetchers into a map that can be easily accessed
	fetchers := map[string]fetcher.Fetcher{
		models.SonarqubeService.String(): sonarqube.NewFetcher(),
	}

	// create the Terraform worker which will do the deployments
	tfw := terraform.NewWorker(s.config.TerraformExecPath)

	// find the template files that will be used to generate the
	// actual deployment files
	templatesDir, err := utils.BuildTemplatesFolderPath()
	if err != nil {
		errors = append(errors, fmt.Errorf("could not build templates directory path: %w", err))
		return nil, nil, errors
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
		filepaths := utils.ReadDir(serviceTemplateDir)

		if len(filepaths) == 0 {
			log.Println(fmt.Errorf("no files founds in dir %s", serviceTemplateDir))
			continue
		}
		templates[dir.Name()] = filepaths
	}

	// if for some reason we failed to find the files
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
			canonical := utils.BuildDeploymentCanonical(nd.UseCase, ns.Service)
			// was previously validated
			st, _ := models.TypeString(ns.Service)

			// We persist the deployment with no instances since we
			// don't know their URLs yet
			d, warning := s.persistDeployment(ctx, canonical, st, nds.CallbackURL)
			if warning != nil {
				warnings = append(warnings, warning)
				continue
			}

			// We create the DTO with the count number of instances
			// on a Pending state
			cd := &dtos.CreatedDeployment{
				Canonical:   canonical,
				Type:        d.Type.String(),
				Instances:   make([]dtos.CreatedInstance, ns.Count),
				CallbackURL: nds.CallbackURL,
			}

			for i := 0; i < ns.Count; i++ {
				cd.Instances[i] = dtos.CreatedInstance{
					State: models.Pending.String(),
				}
			}

			f, ok := fetchers[ns.Service]
			if !ok {
				errors = append(errors, fmt.Errorf("no fetcher found for service %s", ns.Service))
				continue
			}

			key := fmt.Sprintf("%s:%s", nd.UseCase, ns.Service)
			config := configurations[key]

			// fetch the application version the configuration mentions
			appVersion, err := f.GetApplicationVersion(config.Version)
			if err != nil {
				errors = append(errors, err)
				continue
			}

			// fetch the plugins the configuration mentions
			pluginURLs := make([]string, 0)
			for _, p := range config.Plugins {
				plugin, err := f.GetPlugin(p.Name, p.Version)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				pluginURLs = append(pluginURLs, plugin.DownloadURL)
			}

			deployments = append(deployments, cd)

			adminPassword := utils.GeneratePassword(int64(len(canonical)))
			databasePassword := utils.GeneratePassword(int64(ns.Count))

			intpl := utils.Interpolator{
				Name:        canonical,
				Count:       ns.Count,
				DownloadURL: appVersion.DownloadURL,
				Version:     appVersion.Version,
				AdminPass:   adminPassword,
				DbPass:      databasePassword,
				PluginURLs:  pluginURLs,
			}

			// we start goroutines that will do the actual deployments
			// the deployment consists of:
			// 1. Read the template files
			// 2. Replace what needs to be replaced in the template file
			// 3. Persist the resulting file
			// 4. Run the terraform using the terraform worker with the resulting files
			wg.Add(1)
			go func(ns dtos.NewService) {
				defer wg.Done()
				dctx := context.Background()

				// we retrieve the path to the deployment folder
				pathToDir, err := utils.BuildDeploymentFolderPath(cd.Canonical)
				if err != nil {
					errors = append(errors, fmt.Errorf("could not build deployments path: %w", err))
					return
				}

				// we create the directories for the deployment files
				err = os.MkdirAll(pathToDir, 0755)
				if err != nil {
					errors = append(errors, fmt.Errorf("error creating directory for deployment %s: %w", cd.Canonical, err))
					return
				}

				// we retrieve all the filepaths for the files that need
				// interpolation and need to be created on the deployment files folder
				filepaths, ok := templates[ns.Service]
				if !ok {
					errors = append(errors, fmt.Errorf("no template founds for the service %s", ns.Service))
					return
				}

				for _, fp := range filepaths {
					templateFilepath := path.Join(templatesDir, ns.Service, fp)
					f, err := os.Open(templateFilepath)
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

					fcontents = intpl.Interpolate(fcontents)

					deploymentFilepath := path.Join(pathToDir, fp)

					// we take everything that is not the name of the file that will be created and
					// create all the directories needed
					dirsUntilBase := deploymentFilepath[0 : len(deploymentFilepath)-len(path.Base(deploymentFilepath))]
					err = os.MkdirAll(dirsUntilBase, 0755)
					if err != nil {
						log.Fatal(fmt.Errorf("error creating directory for deployment %s: %w", canonical, err))
						return
					}

					// we create the file at the deployment files folder
					cf, err := os.Create(deploymentFilepath)
					if err != nil {
						errors = append(errors, fmt.Errorf("could not create file %s: %w", deploymentFilepath, err))
						return
					}
					defer cf.Close()

					// we write the interpolated contents to the file
					_, err = cf.Write(fcontents)
					if err != nil {
						errors = append(errors, fmt.Errorf("could not write to file %s: %w", deploymentFilepath, err))
						return
					}
				}

				// once all files are interpolated and created we do the
				// deployment logic using the terraform worker
				err = tfw.Deploy(pathToDir)
				if err != nil {
					errors = append(errors, fmt.Errorf("error executing deployment %s: %w", cd.Canonical, err))
					return
				}

				// with the deployment done we fetch the outputs
				// which contain the Public IPs of the deployments
				deploymentURLs, err := tfw.GetIPs(pathToDir)
				if err != nil {
					errors = append(errors, fmt.Errorf("error retrieving the public IPs of the deployment %s: %w", cd.Canonical, err))
					return
				}

				// create the instances with their public IP and a Pending state
				instances := make([]db.Instance, len(deploymentURLs))
				for i, durl := range deploymentURLs {
					dbi := db.Instance{
						DeploymentCanonical: d.Canonical,
						State:               uint(models.Pending),
						URL:                 durl,
						AdminCredential: db.Credential{
							Username: "admin",
							Password: adminPassword,
						},
					}
					instances[i] = dbi
				}

				// persist the new instances
				err = s.repo.BatchCreateInstances(dctx, instances)
				if err != nil {
					errors = append(errors, fmt.Errorf("error updating deployment %s's URL: %w", cd.Canonical, err))
					return
				}
			}(ns)
		}
	}

	// Wait for all the deployment tasks to be over
	// and check for errors in order for them to be logged
	go func() {
		wg.Wait()
		for _, err := range errors {
			log.Println(err)
		}
		for _, w := range warnings {
			log.Println(w)
		}
	}()
	return deployments, warnings, nil
}

func (s *Service) persistDeployment(ctx context.Context, canonical string, st models.Type, callbackURL string) (*models.Deployment, error) {
	dbd, err := s.repo.CreateDeployment(ctx, db.Deployment{
		Canonical:   canonical,
		Instances:   make([]db.Instance, 0),
		Type:        uint(st),
		CallbackURL: callbackURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to persist the deployment %s: %w", canonical, err)
	}
	deployment := &models.Deployment{}
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

// UpdateInstance updates an instance that belongs to a given deployment canonical
// and has a specific instance URL
func (s *Service) UpdateInstance(ctx context.Context, dcan string, i *models.Instance) error {
	err := s.repo.UpdateInstance(ctx, *i.ToDB(dcan))
	if err != nil {
		return fmt.Errorf("failed to update the instance living on the URL %s belonging to the deployment %s: %w", i.URL, dcan, err)
	}
	return nil
}
