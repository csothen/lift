package services

import "github.com/csothen/tmdei-project/internal/db"

type Service struct {
	repo db.Querier
}

func New(repo db.Querier) *Service {
	return &Service{repo}
}
