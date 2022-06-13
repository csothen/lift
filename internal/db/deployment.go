package db

import (
	"context"
	"fmt"
)

func (q *querier) GetAllDeployments(ctx context.Context) ([]*Deployment, error) {
	db := q.db.WithContext(ctx)

	var deployments []*Deployment
	res := db.Preload("Instances").Find(&deployments)
	if res.Error != nil {
		return nil, fmt.Errorf("deployments not found: %w", res.Error)
	}
	return deployments, nil
}

func (q *querier) GetDeploymentByCanonical(ctx context.Context, can string) (*Deployment, error) {
	db := q.db.WithContext(ctx)

	var deployment Deployment
	res := db.Preload("Instances").First(&deployment, "canonical = ?", can)
	if res.Error != nil {
		return nil, fmt.Errorf("deployment not found: %w", res.Error)
	}
	return &deployment, nil
}

func (q *querier) CreateDeployment(ctx context.Context, newD Deployment) (*Deployment, error) {
	db := q.db.WithContext(ctx)

	cres := db.Create(&newD)
	if cres.Error != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", cres.Error)
	}

	var deployment Deployment
	fres := db.Preload("Instances").First(&deployment, "canonical = ?", newD.Canonical)
	if fres.Error != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", fres.Error)
	}
	return &deployment, nil
}

func (q *querier) UpdateDeployment(ctx context.Context, updatedD Deployment) error {
	db := q.db.WithContext(ctx)

	var foundD Deployment
	fres := db.First(&foundD, "canonical = ?", updatedD.Canonical)
	if fres.Error != nil {
		return fmt.Errorf("could not update deployment: %w", fres.Error)
	}

	foundD.Instances = updatedD.Instances
	foundD.CallbackURL = updatedD.CallbackURL

	db.Save(&foundD)
	return nil
}

func (q *querier) DeleteDeployment(ctx context.Context, can string) error {
	db := q.db.WithContext(ctx)

	res := db.Delete(&Deployment{}, "canonical = ?", can)
	if res.Error != nil {
		return fmt.Errorf("failed to delete deployment: %w", res.Error)
	}
	return nil
}
