package models

import "github.com/csothen/lift/internal/db"

type Configuration struct {
	UseCases []UseCaseConfiguration `json:"usecases"`
}

type UseCaseConfiguration struct {
	Name     string                 `json:"name"`
	Services []ServiceConfiguration `json:"services"`
}

type ServiceConfiguration struct {
	Type    Type                `json:"type"`
	Version string              `json:"version"`
	Plugins []PluginInformation `json:"plugins"`
}

type PluginInformation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Configuration) FromDB(dbc *db.Configuration) {
	c.UseCases = make([]UseCaseConfiguration, len(dbc.UseCases))
	for i, dbuc := range dbc.UseCases {
		uc := &UseCaseConfiguration{}
		uc.FromDB(&dbuc)
		c.UseCases[i] = *uc
	}
}

func (c *Configuration) ToDB() *db.Configuration {
	dbc := &db.Configuration{
		UseCases: make([]db.UseCaseConfiguration, len(c.UseCases)),
	}
	for i, uc := range c.UseCases {
		dbc.UseCases[i] = *uc.ToDB()
	}
	return dbc
}

func (uc *UseCaseConfiguration) FromDB(dbuc *db.UseCaseConfiguration) {
	uc.Name = dbuc.Name
	uc.Services = make([]ServiceConfiguration, len(dbuc.Services))
	for i, dbs := range dbuc.Services {
		s := &ServiceConfiguration{}
		s.FromDB(&dbs)
		uc.Services[i] = *s
	}
}

func (uc *UseCaseConfiguration) ToDB() *db.UseCaseConfiguration {
	dbuc := &db.UseCaseConfiguration{
		Name:     uc.Name,
		Services: make([]db.ServiceConfiguration, len(uc.Services)),
	}

	for i, s := range uc.Services {
		dbuc.Services[i] = *s.ToDB(uc.Name)
	}
	return dbuc
}

func (s *ServiceConfiguration) FromDB(dbs *db.ServiceConfiguration) {
	s.Type = Type(dbs.Type)
	s.Version = dbs.Version
	s.Plugins = make([]PluginInformation, len(dbs.Plugins))
	for i, dbp := range dbs.Plugins {
		p := &PluginInformation{}
		p.FromDB(&dbp)
		s.Plugins[i] = *p
	}
}

func (s *ServiceConfiguration) ToDB(usecase string) *db.ServiceConfiguration {
	dbs := &db.ServiceConfiguration{
		UseCase: usecase,
		Type:    uint(s.Type),
		Version: s.Version,
		Plugins: make([]db.PluginInformation, len(s.Plugins)),
	}
	for i, p := range s.Plugins {
		dbs.Plugins[i] = *p.ToDB()
	}
	return dbs
}

func (p *PluginInformation) FromDB(dbp *db.PluginInformation) {
	p.Name = dbp.Name
	p.Version = dbp.Version
}

func (p *PluginInformation) ToDB() *db.PluginInformation {
	return &db.PluginInformation{
		Name:    p.Name,
		Version: p.Version,
	}
}
