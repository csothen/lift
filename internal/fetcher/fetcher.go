package plugins

import "github.com/csothen/tmdei-project/internal/fetcher/types"

type Fetcher interface {
	// Reload will reload all the necessary information from the Internet
	Reload() error

	// GetPlugin takes a name of the plugin and the version wanted and returns
	// either the plugin found or an error in case it does not exist
	GetPlugin(name, version string) (*types.Plugin, error)

	// ListPlugins will list all the plugins available
	ListPlugins() []*types.Plugin

	// GetApplicationVersion will return the requested version of
	// the application or an error in case it does not exist
	GetApplicationVersion(version string) (*types.AppVersion, error)

	// ListApplicationVersions will list all available versions of the
	// application that can be downloaded
	ListApplicationVersions() []*types.AppVersion
}
