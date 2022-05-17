package graph

import (
	"github.com/csothen/tmdei-project/internal/models"
	"github.com/csothen/tmdei-project/internal/models/dtos"
)

func (nuc *NewUseCaseConfiguration) toDTO() *dtos.NewUseCaseConfiguration {
	dto := &dtos.NewUseCaseConfiguration{
		Name:     nuc.Name,
		Services: make([]dtos.NewServiceConfiguration, len(nuc.Services)),
	}

	for i, ns := range nuc.Services {
		sdto := ns.toDTO()
		dto.Services[i] = *sdto
	}
	return dto
}

func (ns *NewServiceConfiguration) toDTO() *dtos.NewServiceConfiguration {
	dto := &dtos.NewServiceConfiguration{
		Type:    ns.Type.String(),
		Version: ns.Version,
		Plugins: make([]dtos.PluginInformation, len(ns.Plugins)),
	}

	for i, pi := range ns.Plugins {
		dto.Plugins[i] = dtos.PluginInformation{
			Name:    pi.Name,
			Version: pi.Version,
		}
	}
	return dto
}

func (nds *NewDeployments) toDTO() *dtos.NewDeployments {
	dto := &dtos.NewDeployments{
		Deployments: make([]dtos.NewDeployment, len(nds.Deployments)),
		CallbackURL: nds.CallbackURL,
	}

	for i, nd := range nds.Deployments {
		ddto := nd.toDTO()
		dto.Deployments[i] = *ddto
	}
	return dto
}

func (nd *NewDeployment) toDTO() *dtos.NewDeployment {
	dto := &dtos.NewDeployment{
		UseCase:  nd.Usecase,
		Services: make([]dtos.NewService, len(nd.Services)),
	}

	for i, ns := range nd.Services {
		dto.Services[i] = dtos.NewService{
			Service: ns.Service,
			Count:   *ns.Count,
		}
	}
	return dto
}

func (c *Configuration) fromModel(cm models.Configuration) {
	c.Usecases = make([]UseCaseConfiguration, len(cm.UseCases))
	for i, ucm := range cm.UseCases {
		var uc *UseCaseConfiguration
		uc.fromModel(ucm)
		c.Usecases[i] = *uc
	}
}

func (uc *UseCaseConfiguration) fromModel(ucm models.UseCaseConfiguration) {
	uc.Name = ucm.Name
	uc.Services = make([]ServiceConfiguration, len(ucm.Services))
	for i, sm := range ucm.Services {
		var s *ServiceConfiguration
		s.fromModel(sm)
		uc.Services[i] = *s
	}
}

func (s *ServiceConfiguration) fromModel(sm models.ServiceConfiguration) {
	s.Type = ServiceType(sm.Type.String())
	s.Version = sm.Version
	s.Plugins = make([]PluginConfiguration, len(sm.Plugins))
	for i, pm := range sm.Plugins {
		s.Plugins[i] = PluginConfiguration{
			Name:    pm.Name,
			Version: pm.Version,
		}
	}
}

func (d *Deployment) fromModel(dm models.Deployment) {
	d.Canonical = dm.Canonical
	d.State = dm.State.String()
	d.Type = dm.Type.String()
	d.URL = dm.URL
	d.CallbackURL = dm.CallbackURL

	var ucred *Credential
	ucred.fromModel(dm.UserCredential)
	d.UserCredential = ucred
}

func (c *Credential) fromModel(cm models.Credential) {
	c.Username = cm.Username
	c.Password = cm.Password
	c.AccessToken = cm.AccessToken
}
