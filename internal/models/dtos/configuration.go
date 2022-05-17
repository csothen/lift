package dtos

type NewUseCaseConfiguration struct {
	Name     string
	Services []NewServiceConfiguration
}

type NewServiceConfiguration struct {
	Type    string
	Version string
	Plugins []PluginInformation
}

type PluginInformation struct {
	Name    string
	Version string
}
