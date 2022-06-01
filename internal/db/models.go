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

	Canonical       string `gorm:"index,unique"`
	State           uint
	Type            uint
	URL             string
	AdminCredential Credential
	UserCredential  Credential
	CallbackURL     string
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

	Name     string                 `gorm:"index,unique"`
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
