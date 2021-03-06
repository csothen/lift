package services

import (
	"context"
	"fmt"

	"github.com/csothen/lift/internal/models"
	"github.com/csothen/lift/internal/models/dtos"
)

// ReadConfiguration reads the whole configuration
func (s *Service) ReadConfiguration(ctx context.Context) *models.Configuration {
	cfg := &models.Configuration{}
	dbConfig, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return cfg
	}
	cfg.FromDB(dbConfig)
	return cfg
}

// ReadConfigurationUseCase reads the information of a specific usecase of the configuration
func (s *Service) ReadConfigurationUseCase(ctx context.Context, usecase string) (*models.UseCaseConfiguration, error) {
	ucCfg := &models.UseCaseConfiguration{}
	dbUcConfig, err := s.repo.GetUseCaseConfiguration(ctx, usecase)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve use case '%s' configuration: %w", usecase, err)
	}
	ucCfg.FromDB(dbUcConfig)
	return ucCfg, nil
}

// ReadConfigurationService reads the information of a specific service in a given usecase of the configuration
func (s *Service) ReadConfigurationService(ctx context.Context, usecase, service string) (*models.ServiceConfiguration, error) {
	sCfg := &models.ServiceConfiguration{}
	stype, err := models.TypeString(service)
	if err != nil {
		return nil, fmt.Errorf("invalid service type %s: %w", service, err)
	}

	dbSConfig, err := s.repo.GetServiceConfiguration(ctx, usecase, uint(stype))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve service configuration for usecase '%s' and service %s: %w", usecase, service, err)
	}
	sCfg.FromDB(dbSConfig)
	return sCfg, nil
}

// AddConfigurationUseCase adds a new usecase to the configuration
func (s *Service) AddConfigurationUseCase(ctx context.Context, ucconfig *dtos.NewUseCaseConfiguration) (*models.Configuration, error) {
	_, err := s.repo.CreateUseCaseConfiguration(ctx, *ucconfig.ToDB())
	if err != nil {
		return nil, fmt.Errorf("failed to add the use case %s to the configuration: %w", ucconfig.Name, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}

// AddConfigurationService adds a new service to a specific usecase in the configuration
func (s *Service) AddConfigurationService(ctx context.Context, usecase string, sconfig *dtos.NewServiceConfiguration) (*models.Configuration, error) {
	_, err := s.repo.CreateServiceConfiguration(ctx, *sconfig.ToDB(usecase))
	if err != nil {
		return nil, fmt.Errorf("failed to add the service %s to the usecase %s in the configuration: %w", sconfig.Type, usecase, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}

// UpdateConfigurationUseCase updates a specific usecase in the configuration
func (s *Service) UpdateConfigurationUsecase(ctx context.Context, usecase string, ucconfig *models.UseCaseConfiguration) (*models.Configuration, error) {
	ucconfig.Name = usecase
	err := s.repo.UpdateUseCaseConfiguration(ctx, *ucconfig.ToDB())
	if err != nil {
		return nil, fmt.Errorf("failed to update the usecase %s in the configuration: %w", usecase, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}

// UpdateConfigurationService updates a service of a specific usecase in the configuration
func (s *Service) UpdateConfigurationService(ctx context.Context, usecase, service string, sconfig *models.ServiceConfiguration) (*models.Configuration, error) {
	st, err := models.TypeString(service)
	if err != nil {
		return nil, fmt.Errorf("invalid service type %s: %w", service, err)
	}

	sconfig.Type = st
	err = s.repo.UpdateServiceConfiguration(ctx, *sconfig.ToDB(usecase))
	if err != nil {
		return nil, fmt.Errorf("failed to update the service %s in the usecase %s in the configuration: %w", service, usecase, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}

// DeleteConfigurationUseCase deletes a specific usecase from the configuration
func (s *Service) DeleteConfigurationUseCase(ctx context.Context, usecase string) (*models.Configuration, error) {
	err := s.repo.DeleteUseCaseConfiguration(ctx, usecase)
	if err != nil {
		return nil, fmt.Errorf("failed to delete usecase %s from the configuration: %w", usecase, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}

// DeleteConfigurationService deletes a service from a specific usecase in the configuration
func (s *Service) DeleteConfigurationService(ctx context.Context, usecase, service string) (*models.Configuration, error) {
	st, err := models.TypeString(service)
	if err != nil {
		return nil, fmt.Errorf("invalid service type %s: %w", service, err)
	}

	err = s.repo.DeleteServiceConfiguration(ctx, usecase, uint(st))
	if err != nil {
		return nil, fmt.Errorf("failed to delete usecase %s from the configuration: %w", usecase, err)
	}

	// fetch updated global configuration
	dbgc, err := s.repo.GetConfiguration(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the global configuration: %w", err)
	}
	gconfig := &models.Configuration{}
	gconfig.FromDB(dbgc)
	return gconfig, nil
}
