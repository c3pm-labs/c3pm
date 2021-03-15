package defaultadapter

import (
	"gopkg.in/yaml.v2"
)

// Config holds the config for the builtin adapter
type Config struct {
	// Sources lists the project's source files.
	Sources []string `yaml:"sources"`
	// Headers lists the project's header files.
	Headers []string `yaml:"headers"`
	// IncludeDirs lists the projects additional header directories.
	IncludeDirs []string `yaml:"include_dirs"`
	// LinuxConfig holds Linux specific configuration
	LinuxConfig *LinuxConfig `yaml:"linux,omitempty"`
}

// LinuxConfig holds specific configuration on Linux operating systems.
type LinuxConfig struct {
	UsePthread bool `yaml:"use_pthread"`
}

// Parse parses a Config
func Parse(c interface{}) (*Config, error) {
	// TODO: is there a better way to do that?
	out, err := yaml.Marshal(c)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(out, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
