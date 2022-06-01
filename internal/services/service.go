package services

import (
	"github.com/csothen/tmdei-project/internal/config"
	"github.com/csothen/tmdei-project/internal/db"
)

type Service struct {
	config *config.Config
	repo   db.Querier
}

func New(repo db.Querier, cfg *config.Config) *Service {
	return &Service{
		repo:   repo,
		config: cfg,
	}
}
