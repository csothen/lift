package services

import (
	"context"
	"fmt"

	"github.com/csothen/tmdei-project/internal/db"
	"github.com/csothen/tmdei-project/internal/helpers"
	"github.com/csothen/tmdei-project/internal/models"
	"github.com/csothen/tmdei-project/internal/models/dtos"
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

	// TODO: Do the validations
	if len(errors) > 0 {
		return nil, nil, errors
	}

	// TODO: For each individual type of deployment start a goroutine that will do the deployment logic

	deployments = make([]*dtos.CreatedDeployment, 0)
	for _, nd := range nds.Deployments {
		for _, ns := range nd.Services {
			for i := 1; i <= ns.Count; i++ {
				canonical := helpers.BuildDeploymentCanonical(nd.UseCase, ns.Service, i)
				// TODO: Send an actual payload to the DB
				dbd, err := s.repo.CreateDeployment(ctx, db.CreateDeploymentParams{})
				if err != nil {
					warnings = append(warnings, fmt.Errorf("failed to persist the deployment %s: %w", canonical, err))
					continue
				}
				var deployment *models.Deployment
				deployment.FromDB(dbd)

				cd := &dtos.CreatedDeployment{
					Canonical: canonical,
					State:     deployment.State.String(),
					Type:      deployment.Type.String(),
				}
				deployments = append(deployments, cd)
			}
		}
	}
	// TODO: Start a goroutine to listen for changes on the deployments
	return deployments, warnings, nil
}

// Update updates a deployment with the given canonical
func (s *Service) UpdateDeployment(ctx context.Context, dcan string, d *models.Deployment) error {
	// TODO: Send an actual payload to the DB
	err := s.repo.UpdateDeployment(ctx, db.UpdateDeploymentParams{})
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
