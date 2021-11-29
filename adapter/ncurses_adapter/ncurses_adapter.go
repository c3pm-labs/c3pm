package ncurses_adapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
)

type Adapter struct {
}

type Config struct {
	Targets []string          `yaml:"targets"`
	Flags   map[string]string `yaml:"variables"`
}

func New() *Adapter {
	return &Adapter{}
}

func executeCli(command string, dir string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("%s process failed: %w", command, err)
	}
	return nil
}

func (a *Adapter) Build(pc *config.ProjectConfig) error {
	err := executeCli("chmod", pc.ProjectRoot, "+x", pc.ProjectRoot+"/configure")
	if err != nil {
		return err
	}
	err = executeCli(
		pc.ProjectRoot+"/configure",
		pc.ProjectRoot,
		"--prefix", pc.ProjectRoot+"/library",
		"--datadir", pc.ProjectRoot+"/library",
	)
	if err != nil {
		return err
	}
	err = executeCli("make", pc.ProjectRoot, "install")
	if err != nil {
		return err
	}
	err = os.Rename(pc.ProjectRoot+"/library/lib/libncurses.a", pc.ProjectRoot+"/libncurses.a")
	if err != nil {
		return err
	}
	return nil
}

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

func (a *Adapter) Targets(pc *config.ProjectConfig) (targets []string, err error) {
	cfg, err := parseConfig(pc.Manifest.Build.Config)
	if err != nil {
		return nil, err
	}

	return cfg.Targets, nil
}

func (a *Adapter) CmakeConfig(pc *config.ProjectConfig) (string, error) {
	return "", nil
}
