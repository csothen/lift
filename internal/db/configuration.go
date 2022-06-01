package db

import (
	"context"
	"fmt"
)

func (q *querier) GetConfiguration(ctx context.Context) (*Configuration, error) {
	var usecases []UseCaseConfiguration
	res := q.db.Find(&usecases)
	if res.Error != nil {
		return nil, fmt.Errorf("configuration not found: %w", res.Error)
	}

	return &Configuration{
		UseCases: usecases,
	}, nil
}

func (q *querier) GetUseCaseConfiguration(ctx context.Context, uc string) (*UseCaseConfiguration, error) {
	var ucconfig UseCaseConfiguration
	res := q.db.First(&ucconfig, "name = ?", uc)
	if res.Error != nil {
		return nil, fmt.Errorf("use case not found: %w", res.Error)
	}
	return &ucconfig, nil
}

func (q *querier) GetServiceConfiguration(ctx context.Context, uc string, service uint) (*ServiceConfiguration, error) {
	var sconfig ServiceConfiguration
	res := q.db.First(&sconfig, "type = ? AND use_case = ?", service, uc)
	if res.Error != nil {
		return nil, fmt.Errorf("service configuration not found: %w", res.Error)
	}
	return &sconfig, nil
}

func (q *querier) CreateUseCaseConfiguration(ctx context.Context, newUC UseCaseConfiguration) (*UseCaseConfiguration, error) {
	cres := q.db.Create(&newUC)
	if cres.Error != nil {
		return nil, fmt.Errorf("failed to create usecase configuration: %w", cres.Error)
	}

	var ucconfig UseCaseConfiguration
	fres := q.db.First(&ucconfig, "name = ?", newUC.Name)
	if fres.Error != nil {
		return nil, fmt.Errorf("failed to create usecase configuration: %w", fres.Error)
	}
	return &ucconfig, nil
}

func (q *querier) CreateServiceConfiguration(ctx context.Context, newS ServiceConfiguration) (*ServiceConfiguration, error) {
	cres := q.db.Create(&newS)
	if cres.Error != nil {
		return nil, fmt.Errorf("failed to create service configuration: %w", cres.Error)
	}

	var sconfig ServiceConfiguration
	fres := q.db.First(&sconfig, "type = ? AND use_case = ?", newS.Type, newS.UseCase)
	if fres.Error != nil {
		return nil, fmt.Errorf("failed to create service configuration: %w", fres.Error)
	}
	return &sconfig, nil
}

func (q *querier) UpdateConfiguration(ctx context.Context, updatedC Configuration) error {
	for _, uc := range updatedC.UseCases {
		var foundUC UseCaseConfiguration
		fres := q.db.First(&foundUC, "name = ?", uc.Name)
		if fres.Error != nil {
			cres := q.db.Create(&uc)
			if cres.Error != nil {
				return fmt.Errorf("could not update configuration: %w", cres.Error)
			}
			continue
		}

		foundUC.Services = uc.Services
		q.db.Save(&foundUC)
	}
	return nil
}

func (q *querier) UpdateUseCaseConfiguration(ctx context.Context, updatedUC UseCaseConfiguration) error {
	var foundUC UseCaseConfiguration
	fres := q.db.First(&foundUC, "name = ?", updatedUC.Name)
	if fres.Error != nil {
		return fmt.Errorf("could not update usecase configuration: %w", fres.Error)
	}

	foundUC.Services = updatedUC.Services
	q.db.Save(&foundUC)
	return nil
}

func (q *querier) UpdateServiceConfiguration(ctx context.Context, updatedS ServiceConfiguration) error {
	var foundS ServiceConfiguration
	fres := q.db.First(&foundS, "type = ? AND use_case = ?", updatedS.Type, updatedS.UseCase)
	if fres.Error != nil {
		return fmt.Errorf("could not update service configuration: %w", fres.Error)
	}

	foundS.Version = updatedS.Version
	foundS.Plugins = updatedS.Plugins
	q.db.Save(&foundS)
	return nil
}

func (q *querier) DeleteUseCaseConfiguration(ctx context.Context, uc string) error {
	res := q.db.Delete(&UseCaseConfiguration{}, "name = ?", uc)
	if res.Error != nil {
		return fmt.Errorf("failed to delete usecase configuration: %w", res.Error)
	}
	return nil
}

func (q *querier) DeleteServiceConfiguration(ctx context.Context, uc string, service uint) error {
	res := q.db.Delete(&ServiceConfiguration{}, "type = ? AND use_case = ?", uc)
	if res.Error != nil {
		return fmt.Errorf("failed to delete service configuration: %w", res.Error)
	}
	return nil
}
