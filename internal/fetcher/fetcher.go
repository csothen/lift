package fetcher

import (
	"github.com/csothen/lift/internal/fetcher/jenkins"
	"github.com/csothen/lift/internal/fetcher/sonarqube"
	"github.com/csothen/lift/internal/fetcher/types"
	"github.com/csothen/lift/internal/models"
)

var (
	Fetchers map[string]Fetcher = map[string]Fetcher{
		models.SonarqubeService.String(): sonarqube.NewFetcher(),
		models.JenkinsService.String():   jenkins.NewFetcher(),
	}
)

type Fetcher interface {
	// Reload will load the data from a static file, if the file is not found
	// it will call the Fetch method to populate the fetcher and create a file
	Reload() error

	// Fetch will fetch all the necessary information from the Internet
	// and persist it in a static file
	Fetch() error

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
