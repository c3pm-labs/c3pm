package cmake_adapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/adapter_interface"
	"github.com/c3pm-labs/c3pm/cmake"
	"github.com/c3pm-labs/c3pm/config"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Adapter struct{}

type Config struct {
	Targets   []string          `yaml:"targets"`
	Variables map[string]string `yaml:"variables"`
}

func (a Adapter) Build(pc *config.ProjectConfig) error {
	cfg, err := parseConfig(pc.Manifest.Build.Config)
	if err != nil {
		return err
	}

	buildDir := filepath.Join(pc.LocalC3PMDirPath(), "build")

	err = cmake.GenerateBuildFiles(pc.ProjectRoot, buildDir, cfg.Variables)
	if err != nil {
		return err
	}

	for _, target := range cfg.Targets {
		err = cmake.Build(buildDir, target)
		if err != nil {
			return err
		}

		err := os.Rename(filepath.Join(buildDir, target), filepath.Join(pc.ProjectRoot, target))
		if err != nil {
			return fmt.Errorf("failed to move target %s to project directory: %v", target, err)
		}
	}

	return nil
}

func (a Adapter) Targets(pc *config.ProjectConfig) ([]string, error) {
	cfg, err := parseConfig(pc.Manifest.Build.Config)
	if err != nil {
		return nil, err
	}

	return cfg.Targets, nil
}

func (a Adapter) CmakeConfig(*config.ProjectConfig) (string, error) {
	return "", nil
}

func New() *Adapter {
	return &Adapter{}
}

var _ adapter_interface.Adapter = (*Adapter)(nil)

func parseConfig(c interface{}) (*Config, error) {
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
