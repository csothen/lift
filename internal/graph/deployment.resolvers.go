package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
)

func (r *mutationResolver) CreateDeployments(ctx context.Context, input NewDeployments) ([]Deployment, error) {
	deployments, _, errors := r.s.CreateDeployment(ctx, input.toDTO())
	if len(errors) > 0 {
		fullErr := fmt.Errorf("failed to create the deployment:")
		for _, err := range errors {
			fullErr = fmt.Errorf("%s\n%w", fullErr, err)
		}
		return nil, fullErr
	}

	nds := make([]Deployment, len(deployments))
	for i, deployment := range deployments {
		nds[i] = Deployment{
			Canonical: deployment.Canonical,
			State:     deployment.State,
			Type:      deployment.Type,
		}
	}
	return nds, nil
}

func (r *queryResolver) Deployments(ctx context.Context) ([]Deployment, error) {
	deployments := r.s.ReadAllDeployments(ctx)

	ds := make([]Deployment, len(deployments))
	for i, deployment := range deployments {
		var d *Deployment
		d.fromModel(*deployment)
		ds[i] = *d
	}
	return ds, nil
}

func (r *queryResolver) FindDeployment(ctx context.Context, canonical string) (*Deployment, error) {
	deployment, err := r.s.ReadOneDeployment(ctx, canonical)
	if err != nil {
		return nil, err
	}

	var d *Deployment
	d.fromModel(*deployment)
	return d, nil
}
