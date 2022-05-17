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

	Canonical       string
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
	gorm.Model

	UseCases []UseCaseConfiguration
}

type UseCaseConfiguration struct {
	gorm.Model

	Name     string
	Services []ServiceConfiguration
}

type ServiceConfiguration struct {
	gorm.Model

	Type    uint
	Version string
	Plugins []PluginConfiguration
}

type PluginConfiguration struct {
	gorm.Model

	Name    string
	Version string
}
