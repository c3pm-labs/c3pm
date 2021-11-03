//Package manifest handles loading and writing of the C3PM manifest file (usually c3pm.yml).
//It also stores the various types supported by the manifest, as well as utility functions to use them.
package manifest

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Manifest is the main configuration structure for C3PM.
type Manifest struct {
	C3PMVersion   C3PMVersion    `yaml:"c3pm_version"`
	Type          Type           `yaml:"type"`
	Name          string         `yaml:"name"`
	Description   string         `yaml:"description"`
	Version       Version        `yaml:"version"`
	Publish       *PublishConfig `yaml:"publish,omitempty"`
	Build         *BuildConfig   `yaml:"build,omitempty"`
	Documentation string         `yaml:"documentation"`
	Website       string         `yaml:"website"`
	Repository    string         `yaml:"repository"`
	Contributors  string         `yaml:"contributors"`
	Standard      string         `yaml:"standard"`
	License       string         `yaml:"license"`
	Dependencies  Dependencies   `yaml:"dependencies"`
	Tags          []string       `yaml:"tags,omitempty"`
}

type BuildConfig struct {
	Adapter *AdapterConfig `yaml:"adapter,omitempty"`
	Config  interface{}    `yaml:"config,omitempty"`
}

type AdapterConfig struct {
	Name    string  `yaml:"name"`
	Version Version `yaml:"version,omitempty"`
}

type PublishConfig struct {
	IncludeDirs []string `yaml:"include_dirs,omitempty"`
	Include     []string `yaml:"include,omitempty"`
	Exclude     []string `yaml:"exclude,omitempty"`
}

// New returns the default manifest values.
func New() Manifest {
	c3pmAdapterVersion, _ := VersionFromString("0.0.1")
	defaultManifest := Manifest{
		C3PMVersion:  C3PMVersion1,
		Dependencies: make(map[string]string),
		Standard:     "20",
		Publish: &PublishConfig{
			Exclude: []string{},
			Include: []string{"**/**"},
		},
		Build: &BuildConfig{
			Adapter: &AdapterConfig{
				Name:    "c3pm",
				Version: c3pmAdapterVersion,
			},
			Config: &struct {
				Sources     []string
				Headers     []string
				IncludeDirs []string `yaml:"include_dirs"`
			}{
				Sources:     []string{"**/*.cpp"},
				Headers:     []string{"**/*.hpp"},
				IncludeDirs: []string{"include"},
			},
		},
	}
	return defaultManifest
}

func deserialize(config []byte) (Manifest, error) {
	man := New()
	err := yaml.Unmarshal(config, &man)
	if err != nil {
		return man, err
	}
	return man, nil
}

func checkManifest(man Manifest) error {
	if man.Name == "" {
		return errors.New("Field \"name\" is missing")
	}
	if man.Type != "executable" && man.Type != "library" {
		return errors.New("Field \"type\" is missing or incorrect [executable | library]")
	}
	if man.Version.Version == nil {
		return errors.New("Field \"version\" is missing")
	}
	return nil
}

// Load reads the file located at the given path, and stores its contents in a new Manifest struct.
func Load(path string) (Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	man, err := deserialize(data)
	if err != nil {
		return man, err
	}
	if err = checkManifest(man); err != nil {
		return man, err
	}
	return man, err
}

func (m *Manifest) serialize() ([]byte, error) {
	return yaml.Marshal(&m)
}

// Save writes the current Manifest struct into the destination path.
func (m *Manifest) Save(destination string) error {
	data, err := m.serialize()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destination, data, 0644)
}

// Targets returns the CMake targets to use.
// TODO: delete this
func (m *Manifest) Targets() []string {
	return []string{m.Name}
}
