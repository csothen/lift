package db

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/csothen/lift/internal/config"
)

type Querier interface {
	// API Key operations
	ValidateAPIKey(ctx context.Context, key string) error
	RefreshAPIKey(ctx context.Context, newKey string) error

	// Configuration CRUD operations
	GetConfiguration(ctx context.Context) (*Configuration, error)
	GetUseCaseConfiguration(ctx context.Context, uc string) (*UseCaseConfiguration, error)
	GetServiceConfiguration(ctx context.Context, uc string, service uint) (*ServiceConfiguration, error)
	CreateUseCaseConfiguration(ctx context.Context, newUC UseCaseConfiguration) (*UseCaseConfiguration, error)
	CreateServiceConfiguration(ctx context.Context, newS ServiceConfiguration) (*ServiceConfiguration, error)
	UpdateConfiguration(ctx context.Context, updatedC Configuration) error
	UpdateUseCaseConfiguration(ctx context.Context, updatedUC UseCaseConfiguration) error
	UpdateServiceConfiguration(ctx context.Context, updatedS ServiceConfiguration) error
	DeleteUseCaseConfiguration(ctx context.Context, uc string) error
	DeleteServiceConfiguration(ctx context.Context, uc string, service uint) error

	// Deployment operations
	GetAllDeployments(ctx context.Context) ([]*Deployment, error)
	GetDeploymentByCanonical(ctx context.Context, can string) (*Deployment, error)
	CreateDeployment(ctx context.Context, newD Deployment) (*Deployment, error)
	UpdateDeployment(ctx context.Context, updatedD Deployment) error
	DeleteDeployment(ctx context.Context, can string) error

	// Instances operations
	BatchCreateInstances(ctx context.Context, instances []Instance) error
	UpdateInstance(ctx context.Context, updatedI Instance) error
}

type querier struct {
	db *gorm.DB
}

func New(cfg *config.Config) Querier {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(
		&APIKey{},
		&Deployment{}, &Instance{}, &Credential{},
		&UseCaseConfiguration{}, &ServiceConfiguration{}, &PluginInformation{},
	)

	return &querier{db}
}
