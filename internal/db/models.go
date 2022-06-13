package db

import (
	"gorm.io/gorm"
)

type APIKey struct {
	gorm.Model

	Value string
}

type Deployment struct {
	gorm.Model

	Canonical   string `gorm:"uniqueIndex"`
	Type        uint
	Instances   []Instance `gorm:"foreignKey:DeploymentCanonical;references:Canonical"`
	CallbackURL string
}

type Instance struct {
	gorm.Model

	DeploymentCanonical string
	URL                 string `gorm:"uniqueIndex"`
	State               uint
	AdminCredentialID   *uint
	AdminCredential     Credential
	UserCredentialID    *uint
	UserCredential      Credential
}

type Credential struct {
	gorm.Model

	Username    string
	Password    string
	AccessToken string
}

type Configuration struct {
	UseCases []UseCaseConfiguration
}

type UseCaseConfiguration struct {
	gorm.Model

	Name     string                 `gorm:"uniqueIndex"`
	Services []ServiceConfiguration `gorm:"foreignKey:UseCase;references:Name"`
}

type ServiceConfiguration struct {
	gorm.Model

	UseCase string
	Type    uint
	Version string
	Plugins []PluginInformation `gorm:"many2many:service_configuration_plugins;"`
}

type PluginInformation struct {
	gorm.Model

	Name    string
	Version string
}
