// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"fmt"
	"io"
	"strconv"
)

type Configuration struct {
	Usecases []UseCaseConfiguration `json:"usecases"`
}

type Credential struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string `json:"accessToken"`
}

type Deployment struct {
	Canonical      string      `json:"canonical"`
	State          string      `json:"state"`
	Type           string      `json:"type"`
	URL            string      `json:"url"`
	UserCredential *Credential `json:"userCredential"`
	CallbackURL    string      `json:"callbackURL"`
}

type NewDeployment struct {
	Usecase  string       `json:"usecase"`
	Services []NewService `json:"services"`
}

type NewDeployments struct {
	Deployments []NewDeployment `json:"deployments"`
	CallbackURL string          `json:"callbackURL"`
}

type NewPluginConfiguration struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type NewService struct {
	Service string `json:"service"`
	Count   *int   `json:"count"`
}

type NewServiceConfiguration struct {
	Usecase string                   `json:"usecase"`
	Type    ServiceType              `json:"type"`
	Version string                   `json:"version"`
	Plugins []NewPluginConfiguration `json:"plugins"`
}

type NewUseCaseConfiguration struct {
	Name     string                    `json:"name"`
	Services []NewServiceConfiguration `json:"services"`
}

type PluginConfiguration struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServiceConfiguration struct {
	Type    ServiceType           `json:"type"`
	Version string                `json:"version"`
	Plugins []PluginConfiguration `json:"plugins"`
}

type UseCaseConfiguration struct {
	Name     string                 `json:"name"`
	Services []ServiceConfiguration `json:"services"`
}

type DeploymentState string

const (
	DeploymentStateRunning DeploymentState = "Running"
	DeploymentStatePending DeploymentState = "Pending"
	DeploymentStateStopped DeploymentState = "Stopped"
)

var AllDeploymentState = []DeploymentState{
	DeploymentStateRunning,
	DeploymentStatePending,
	DeploymentStateStopped,
}

func (e DeploymentState) IsValid() bool {
	switch e {
	case DeploymentStateRunning, DeploymentStatePending, DeploymentStateStopped:
		return true
	}
	return false
}

func (e DeploymentState) String() string {
	return string(e)
}

func (e *DeploymentState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeploymentState(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeploymentState", str)
	}
	return nil
}

func (e DeploymentState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ServiceType string

const (
	ServiceTypeSonarqube ServiceType = "Sonarqube"
)

var AllServiceType = []ServiceType{
	ServiceTypeSonarqube,
}

func (e ServiceType) IsValid() bool {
	switch e {
	case ServiceTypeSonarqube:
		return true
	}
	return false
}

func (e ServiceType) String() string {
	return string(e)
}

func (e *ServiceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ServiceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ServiceType", str)
	}
	return nil
}

func (e ServiceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
