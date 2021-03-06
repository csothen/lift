package graph

import (
	"github.com/csothen/lift/internal/models"
	"github.com/csothen/lift/internal/models/dtos"
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

func (uuc *UpdateUseCaseConfiguration) toModel(name string) *models.UseCaseConfiguration {
	ucm := &models.UseCaseConfiguration{
		Name:     name,
		Services: make([]models.ServiceConfiguration, len(uuc.Services)),
	}

	for i, ns := range uuc.Services {
		smodel := *ns.toModel()
		ucm.Services[i] = smodel
	}

	return ucm
}

func (ns *NewServiceConfiguration) toDTO() *dtos.NewServiceConfiguration {
	dto := &dtos.NewServiceConfiguration{
		Type:    ns.Type.String(),
		Plugins: make([]dtos.PluginInformation, len(ns.Plugins)),
	}

	if ns.Version != nil {
		dto.Version = *ns.Version
	}

	for i, pi := range ns.Plugins {
		pidto := dtos.PluginInformation{
			Name: pi.Name,
		}

		if pi.Version != nil {
			pidto.Version = *pi.Version
		}

		dto.Plugins[i] = pidto
	}
	return dto
}

func (ns *NewServiceConfiguration) toModel() *models.ServiceConfiguration {
	st, _ := models.TypeString(ns.Type.String())
	sm := &models.ServiceConfiguration{
		Type:    st,
		Plugins: make([]models.PluginInformation, len(ns.Plugins)),
	}

	if ns.Version != nil {
		sm.Version = *ns.Version
	}

	for i, np := range ns.Plugins {
		npm := models.PluginInformation{
			Name: np.Name,
		}

		if np.Version != nil {
			npm.Version = *np.Version
		}

		sm.Plugins[i] = npm
	}
	return sm
}

func (us *UpdateServiceConfiguration) toModel(service string) *models.ServiceConfiguration {
	st, _ := models.TypeString(service)
	sm := &models.ServiceConfiguration{
		Type:    st,
		Plugins: make([]models.PluginInformation, len(us.Plugins)),
	}

	if us.Version != nil {
		sm.Version = *us.Version
	}

	for i, np := range us.Plugins {
		npm := models.PluginInformation{
			Name: np.Name,
		}

		if np.Version != nil {
			npm.Version = *np.Version
		}

		sm.Plugins[i] = npm
	}
	return sm
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
			Service: ns.Service.String(),
			Count:   *ns.Count,
		}
	}
	return dto
}

func (c *Configuration) fromModel(cm models.Configuration) {
	c.Usecases = make([]UseCaseConfiguration, len(cm.UseCases))
	for i, ucm := range cm.UseCases {
		uc := &UseCaseConfiguration{}
		uc.fromModel(ucm)
		c.Usecases[i] = *uc
	}
}

func (uc *UseCaseConfiguration) fromModel(ucm models.UseCaseConfiguration) {
	uc.Name = ucm.Name
	uc.Services = make([]ServiceConfiguration, len(ucm.Services))
	for i, sm := range ucm.Services {
		s := &ServiceConfiguration{}
		s.fromModel(sm)
		uc.Services[i] = *s
	}
}

func (s *ServiceConfiguration) fromModel(sm models.ServiceConfiguration) {
	s.Type = ServiceType(sm.Type.String())
	s.Plugins = make([]PluginConfiguration, len(sm.Plugins))
	if sm.Version != "" {
		s.Version = &sm.Version
	}

	for i, pm := range sm.Plugins {
		pc := PluginConfiguration{
			Name: pm.Name,
		}

		if pm.Version != "" {
			pc.Version = &pm.Version
		}

		s.Plugins[i] = pc
	}
}

func (d *Deployment) fromModel(dm models.Deployment) {
	d.Canonical = dm.Canonical
	d.Instances = make([]Instance, len(dm.Instances))
	d.Type = dm.Type.String()
	d.CallbackURL = dm.CallbackURL

	for i, im := range dm.Instances {
		instance := &Instance{}
		instance.fromModel(im)
		d.Instances[i] = *instance
	}
}

func (i *Instance) fromModel(im models.Instance) {
	i.URL = im.URL
	i.State = DeploymentState(im.State.String())

	ucred := &Credential{}
	ucred.fromModel(im.UserCredential)
	i.UserCredential = ucred
}

func (c *Credential) fromModel(cm models.Credential) {
	c.Username = cm.Username
	c.Password = cm.Password
}
