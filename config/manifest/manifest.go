package manifest

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type LinuxConfig struct {
	UsePthread bool `yaml:"pthread"`
}

type Manifest struct {
	C3pmVersion  C3pmVersion  `yaml:"c3pm_version"`
	Type         Type         `yaml:"type"`
	Name         string       `yaml:"name"`
	Description  string       `yaml:"description"`
	Version      Version      `yaml:"version"`
	Standard     string       `yaml:"standard"`
	License      string       `yaml:"license"`
	Files        FilesConfig  `yaml:"files"`
	Include      []string     `yaml:"include"`
	Exclude      []string     `yaml:"exclude"`
	Dependencies Dependencies `yaml:"dependencies"`
	CustomCmake  *CustomCmake `yaml:"custom_cmake,omitempty"`
	LinuxConfig  *LinuxConfig `yaml:"linux,omitempty"`
}

var defaultManifest = Manifest{
	C3pmVersion: C3pmVersion1,
	Files: FilesConfig{
		Sources:             []string{"**/*.cpp"},
		Includes:            []string{"**/*.hpp"},
		ExportedIncludeDirs: []string{},
	},
	Dependencies: make(map[string]string),
	Standard:     "20",
}

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

func Load(path string) (Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	m, err := deserialize(data)
	if err != nil {
		return Manifest{}, err
	}

	//if m.CustomCmake != nil && !m.Files.IsEmpty() {
	//	return Manifest{}, fmt.Errorf("cannot specify custom_cmake and source files")
	//}

	return m, nil
}

func (m *Manifest) serialize() ([]byte, error) {
	return yaml.Marshal(&m)
}

func (m *Manifest) Save(destination string) error {
	data, err := m.serialize()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(destination, data, 0644)
}

func (m *Manifest) Targets() []string {
	if m.CustomCmake != nil {
		return m.CustomCmake.Targets
	}
	return []string{m.Name}
}
