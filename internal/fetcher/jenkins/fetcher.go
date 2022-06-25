package jenkins

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/csothen/lift/internal/fetcher/types"
	"github.com/csothen/lift/internal/utils"
)

const (
	StaticContentFilename string = "fetched/jenkins.json"
)

type fetcher struct {
	Plugins         map[string]*types.Plugin     `json:"plugins"`
	JenkinsVersions map[string]*types.AppVersion `json:"jenkins"`
}

func NewFetcher() *fetcher {
	return &fetcher{
		Plugins:         make(map[string]*types.Plugin),
		JenkinsVersions: make(map[string]*types.AppVersion),
	}
}

func (f *fetcher) Reload() error {
	if f.hasStaticData() {
		return f.readData()
	}
	return f.Fetch()
}

func (f *fetcher) Fetch() error {
	return f.writeData()
}

func (f *fetcher) GetPlugin(name, version string) (*types.Plugin, error) {
	return &types.Plugin{}, nil
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
	// TODO: Currently a constant but could be dynamically fetched in the future
	return &types.AppVersion{
		Version: "2.332",
	}, nil
}

func (f *fetcher) ListApplicationVersions() []*types.AppVersion {
	if !f.hasData() {
		if err := f.Reload(); err != nil {
			return nil
		}
	}

	versions := make([]*types.AppVersion, 0)
	for _, v := range f.JenkinsVersions {
		versions = append(versions, v)
	}
	return versions
}

func (f *fetcher) hasData() bool {
	return len(f.Plugins) > 0 && len(f.JenkinsVersions) > 0
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
