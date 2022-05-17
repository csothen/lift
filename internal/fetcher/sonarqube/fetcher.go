package sonarqube

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/csothen/tmdei-project/internal/fetcher/types"
)

type fetcher struct {
	plugins           map[string]*types.Plugin
	sonarqubeVersions map[string]*types.AppVersion
}

func NewFetcher() *fetcher {
	return &fetcher{
		plugins:           make(map[string]*types.Plugin),
		sonarqubeVersions: make(map[string]*types.AppVersion),
	}
}

func (f *fetcher) Reload() error {
	psURL, err := url.Parse(pluginsSource)
	if err != nil {
		return fmt.Errorf("plugins source is not a valid URL: %w", err)
	}

	ssURL, err := url.Parse(sonarqubeSource)
	if err != nil {
		return fmt.Errorf("sonarqube versions source is not a valid URL: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go f.fetchPlugins(&wg, psURL)
	go f.fetchSonarqubeVersions(&wg, ssURL)
	wg.Wait()

	return nil
}

func (f *fetcher) GetPlugin(name, version string) (*types.Plugin, error) {
	plugin, ok := f.plugins[buildPluginKey(name, version)]
	if !ok {
		return nil, fmt.Errorf("plugin with name %s and version %s was not found", name, version)
	}
	return plugin, nil
}

func (f *fetcher) ListPlugins() []*types.Plugin {
	plugins := make([]*types.Plugin, 0)
	for _, v := range f.plugins {
		plugins = append(plugins, v)
	}
	return plugins
}

func (f *fetcher) GetApplicationVersion(version string) (*types.AppVersion, error) {
	appVersion, ok := f.sonarqubeVersions[version]
	if !ok {
		return nil, fmt.Errorf("sonarqube version %s was not found", version)
	}
	return appVersion, nil
}

func (f *fetcher) ListApplicationVersions() []*types.AppVersion {
	versions := make([]*types.AppVersion, 0)
	for _, v := range f.sonarqubeVersions {
		versions = append(versions, v)
	}
	return versions
}
