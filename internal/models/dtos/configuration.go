package dtos

import (
	"github.com/csothen/lift/internal/db"
	"github.com/csothen/lift/internal/models"
)

type NewUseCaseConfiguration struct {
	Name     string
	Services []NewServiceConfiguration
}

func (uc *NewUseCaseConfiguration) ToDB() *db.UseCaseConfiguration {
	dbuc := &db.UseCaseConfiguration{
		Name:     uc.Name,
		Services: make([]db.ServiceConfiguration, len(uc.Services)),
	}
	for i, s := range uc.Services {
		dbuc.Services[i] = *s.ToDB(uc.Name)
	}
	return dbuc
}

type NewServiceConfiguration struct {
	Type    string
	Version string
	Plugins []PluginInformation
}

func (s *NewServiceConfiguration) ToDB(usecase string) *db.ServiceConfiguration {
	st, _ := models.TypeString(s.Type)
	dbs := &db.ServiceConfiguration{
		UseCase: usecase,
		Type:    uint(st),
		Version: s.Version,
		Plugins: make([]db.PluginInformation, len(s.Plugins)),
	}
	for i, p := range s.Plugins {
		dbs.Plugins[i] = *p.ToDB()
	}
	return dbs
}

type PluginInformation struct {
	Name    string
	Version string
}

func (p *PluginInformation) ToDB() *db.PluginInformation {
	return &db.PluginInformation{
		Name:    p.Name,
		Version: p.Version,
	}
}
