package sonarqube

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/csothen/lift/internal/fetcher/types"
	"github.com/csothen/lift/internal/utils"
)

const (
	StaticContentFilename string = "fetched/sonarqube.json"
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
	if f.hasStaticData() {
		return f.readData()
	}
	return f.Fetch()
}

func (f *fetcher) Fetch() error {
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
	if !f.hasData() {
		if err := f.Reload(); err != nil {
			return nil, fmt.Errorf("could not reload plugins data: %w", err)
		}
	}

	plugin, ok := f.Plugins[buildPluginKey(name, version)]
	if !ok {
		return nil, fmt.Errorf("plugin with name %s and version %s was not found", name, version)
	}
	return plugin, nil
}

func (f *fetcher) ListPlugins() []*types.Plugin {
	if !f.hasData() {
		if err := f.Reload(); err != nil {
			return nil
		}
	}

	plugins := make([]*types.Plugin, 0)
	for _, v := range f.Plugins {
		plugins = append(plugins, v)
	}
	return plugins
}

func (f *fetcher) GetApplicationVersion(version string) (*types.AppVersion, error) {
	if !f.hasData() {
		if err := f.Reload(); err != nil {
			return nil, fmt.Errorf("could not reload sonarqube versions data: %w", err)
		}
	}

	appVersion, ok := f.SonarqubeVersions[version]
	if !ok {
		return nil, fmt.Errorf("sonarqube version %s was not found", version)
	}
	return appVersion, nil
}

func (f *fetcher) ListApplicationVersions() []*types.AppVersion {
	if !f.hasData() {
		if err := f.Reload(); err != nil {
			return nil
		}
	}

	versions := make([]*types.AppVersion, 0)
	for _, v := range f.SonarqubeVersions {
		versions = append(versions, v)
	}
	return versions
}

func (f *fetcher) hasData() bool {
	return len(f.Plugins) > 0 && len(f.SonarqubeVersions) > 0
}

func (f *fetcher) hasStaticData() bool {
	staticPath, err := utils.BuildStaticFolderPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path.Join(staticPath, StaticContentFilename))
	return err == nil
}

func (f *fetcher) readData() error {
	staticPath, err := utils.BuildStaticFolderPath()
	if err != nil {
		return fmt.Errorf("could retrieve static folder path: %w", err)
	}

	file, err := os.Open(path.Join(staticPath, StaticContentFilename))
	if err != nil {
		return fmt.Errorf("could not open static content file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(f)
	if err != nil {
		return fmt.Errorf("could not decode file data: %w", err)
	}
	return nil
}

func (f *fetcher) writeData() error {
	staticPath, err := utils.BuildStaticFolderPath()
	if err != nil {
		return fmt.Errorf("could retrieve static folder path: %w", err)
	}

	// Make sure the static folder exists
	err = os.MkdirAll(staticPath, 0755)
	if err != nil {
		return fmt.Errorf("could not create static folder: %w", err)
	}

	file, err := os.Create(path.Join(staticPath, StaticContentFilename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(f)
	if err != nil {
		return fmt.Errorf("could not encode fetcher data: %w", err)
	}
	return nil
}
