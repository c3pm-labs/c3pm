package config

import "errors"

// Config holds the config for the builtin adapter
type Config struct {
	// Sources lists the project's source files.
	Sources []string
	// Headers lists the project's header files.
	Headers []string
	// IncludeDirs lists the projects additional header directories.
	IncludeDirs []string
	// Standard is the c++ standard used to build the project
	Standard     string
	LinuxConfig  *LinuxConfig
}

// LinuxConfig holds specific configuration on Linux operating systems.
type LinuxConfig struct {
	UsePthread bool
}

func Parse(c interface{}) (*Config, error) {
	cfg, ok := c.(*Config)
	if !ok {
		return nil, errors.New("failed to parse build config")
	}
	return cfg, nil
}
