package models

import "github.com/csothen/tmdei-project/internal/db"

type Configuration struct {
	UseCases []UseCaseConfiguration `json:"usecases"`
}

type UseCaseConfiguration struct {
	Name     string                 `json:"name"`
	Services []ServiceConfiguration `json:"services"`
}

type ServiceConfiguration struct {
	Type    Type                  `json:"type"`
	Version string                `json:"version"`
	Plugins []PluginConfiguration `json:"plugins"`
}

type PluginConfiguration struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (c *Configuration) FromDB(dbc db.Configuration) {
	c.UseCases = make([]UseCaseConfiguration, len(dbc.UseCases))
	for _, dbuc := range dbc.UseCases {
		var uc UseCaseConfiguration
		uc.FromDB(dbuc)
		c.UseCases = append(c.UseCases, uc)
	}
}

func (uc *UseCaseConfiguration) FromDB(dbuc db.UseCaseConfiguration) {
	uc.Name = dbuc.Name
	uc.Services = make([]ServiceConfiguration, len(dbuc.Services))
	for _, dbs := range dbuc.Services {
		var s ServiceConfiguration
		s.FromDB(dbs)
		uc.Services = append(uc.Services, s)
	}
}

func (s *ServiceConfiguration) FromDB(dbs db.ServiceConfiguration) {
	s.Type = Type(dbs.Type)
	s.Version = dbs.Version
	s.Plugins = make([]PluginConfiguration, len(dbs.Plugins))
	for _, dbp := range dbs.Plugins {
		var p PluginConfiguration
		p.FromDB(dbp)
		s.Plugins = append(s.Plugins, p)
	}
}

func (p *PluginConfiguration) FromDB(dbp db.PluginConfiguration) {
	p.Name = dbp.Name
	p.Version = dbp.Version
}
