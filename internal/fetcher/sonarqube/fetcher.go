package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/csothen/tmdei-project/internal/fetcher/types"
)

type fetcher struct {
	Plugins           map[string]*types.Plugin     `json:"plugins"`
	SonarqubeVersions map[string]*types.AppVersion `json:"sonarqube"`
}

func NewFetcher() *fetcher {
	return &fetcher{
		Plugins:           make(map[string]*types.Plugin),
		SonarqubeVersions: make(map[string]*types.AppVersion),
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

	return f.writeData()
}

func (f *fetcher) GetPlugin(name, version string) (*types.Plugin, error) {
	plugin, ok := f.Plugins[buildPluginKey(name, version)]
	if !ok {
		return nil, fmt.Errorf("plugin with name %s and version %s was not found", name, version)
	}
	return plugin, nil
}

func (f *fetcher) ListPlugins() []*types.Plugin {
	plugins := make([]*types.Plugin, 0)
	for _, v := range f.Plugins {
		plugins = append(plugins, v)
	}
	return plugins
}

func (f *fetcher) GetApplicationVersion(version string) (*types.AppVersion, error) {
	appVersion, ok := f.SonarqubeVersions[version]
	if !ok {
		return nil, fmt.Errorf("sonarqube version %s was not found", version)
	}
	return appVersion, nil
}

func (f *fetcher) ListApplicationVersions() []*types.AppVersion {
	versions := make([]*types.AppVersion, 0)
	for _, v := range f.SonarqubeVersions {
		versions = append(versions, v)
	}
	return versions
}

func (f *fetcher) writeData() error {
	d, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	staticPath := path.Join(path.Join(d, "static"))
	err = os.MkdirAll(staticPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create static folder: %w", err)
	}

	file, err := os.Create(path.Join(staticPath, "sonarqube.json"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(f)

	return nil
}
