package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) AddUseCaseConfiguration(ctx context.Context, input NewUseCaseConfiguration) (*Configuration, error) {
	configuration, err := r.s.AddConfigurationUseCase(ctx, input.toDTO())
	if err != nil {
		return nil, err
	}

	cfg := &Configuration{}
	cfg.fromModel(*configuration)
	return cfg, nil
}

func (r *mutationResolver) AddServiceConfiguration(ctx context.Context, input NewServiceConfiguration) (*Configuration, error) {
	configuration, err := r.s.AddConfigurationService(ctx, input.Usecase, input.toDTO())
	if err != nil {
		return nil, err
	}

	cfg := &Configuration{}
	cfg.fromModel(*configuration)
	return cfg, nil
}

func (r *queryResolver) Configuration(ctx context.Context) (*Configuration, error) {
	configuration := r.s.ReadConfiguration(ctx)

	cfg := &Configuration{}
	cfg.fromModel(*configuration)
	return cfg, nil
}

func (r *queryResolver) FindUseCaseConfiguration(ctx context.Context, name string) (*UseCaseConfiguration, error) {
	ucConfiguration, err := r.s.ReadConfigurationUseCase(ctx, name)
	if err != nil {
		return nil, err
	}

	ucc := &UseCaseConfiguration{}
	ucc.fromModel(*ucConfiguration)
	return ucc, nil
}

func (r *queryResolver) FindServiceConfiguration(ctx context.Context, uc string, service string) (*ServiceConfiguration, error) {
	sConfiguration, err := r.s.ReadConfigurationService(ctx, uc, service)
	if err != nil {
		return nil, err
	}

	sc := &ServiceConfiguration{}
	sc.fromModel(*sConfiguration)
	return sc, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
