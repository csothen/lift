package db

import (
	"context"
	"fmt"
)

func (q *querier) BatchCreateInstances(ctx context.Context, instances []Instance) error {
	db := q.db.WithContext(ctx)

	cres := db.Create(&instances)
	if cres.Error != nil {
		return fmt.Errorf("failed to create instances: %w", cres.Error)
	}
	return nil
}

func (q *querier) UpdateInstance(ctx context.Context, updatedI Instance) error {
	db := q.db.WithContext(ctx)

	var foundI Instance
	fres := db.First(&foundI, "deployment_canonical = ? AND url = ?", updatedI.DeploymentCanonical, updatedI.URL)
	if fres.Error != nil {
		return fmt.Errorf("could not update deployment: %w", fres.Error)
	}

	foundI.State = updatedI.State
	foundI.AdminCredential = updatedI.AdminCredential
	foundI.UserCredential = updatedI.UserCredential

	db.Save(&foundI)
	return nil
}
