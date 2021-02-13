//Package manifest handles loading and writing of the C3PM manifest file (usually c3pm.yml).
//It also stores the various types supported by the manifest, as well as utility functions to use them.
package manifest

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Manifest is the main configuration structure for C3PM.
type Manifest struct {
	C3PMVersion  C3PMVersion  `yaml:"c3pm_version"`
	Type         Type         `yaml:"type"`
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	Version      Version      `yaml:"version"`
	Standard     string       `yaml:"standard"`
	License      string       `yaml:"license"`
	Files        FilesConfig  `yaml:"files"`
	Dependencies Dependencies `yaml:"dependencies"`
	CustomCMake  *CustomCMake `yaml:"custom_cmake,omitempty"`
	LinuxConfig  *LinuxConfig `yaml:"linux,omitempty"`
}

// LinuxConfig holds specific configuration on Linux operating systems.
type LinuxConfig struct {
	UsePthread bool `yaml:"pthread"`
}

var defaultManifest = Manifest{
	C3PMVersion: C3PMVersion1,
	Files: FilesConfig{
		Sources:             []string{"**/*.cpp"},
		Includes:            []string{"**/*.hpp"},
		ExportedIncludeDirs: []string{},
	},
	Dependencies: make(map[string]string),
	Standard:     "20",
}

// New returns the default manifest values.
func New() Manifest {
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

// Loads reads the file located at the given path, and stores its contents in a new Manifest struct.
func Load(path string) (Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	m, err := deserialize(data)
	if err != nil {
		return Manifest{}, err
	}

	//if m.CustomCMake != nil && !m.Files.IsEmpty() {
	//	return Manifest{}, fmt.Errorf("cannot specify custom_cmake and source files")
	//}

	return m, nil
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
func (m *Manifest) Targets() []string {
	if m.CustomCMake != nil {
		return m.CustomCMake.Targets
	}
	return []string{m.Name}
}
