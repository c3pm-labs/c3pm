package configure_adapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"os"
	"os/exec"
)

type ConfigureAdapter struct {
}

// New creates a new builtin MakefileAdapter
func New() *ConfigureAdapter {
	return &ConfigureAdapter{}
}


func executeCli(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("%s process failed: %w", command, err)
	}
	return nil
}

func (a *ConfigureAdapter) Build(pc *config.ProjectConfig) error {
	err := executeCli("chmod", "+x", pc.ProjectRoot + "/configure")
	if err != nil {
		return err
	}
	err = executeCli(pc.ProjectRoot + "/configure", "--prefix", pc.ProjectRoot, "--exec-prefix", pc.ProjectRoot)
	if err != nil {
		return err
	}
	err = executeCli("make", "-C", pc.ProjectRoot)
	if err != nil {
		return err
	}
	err = os.Rename(pc.ProjectRoot + "/lib/libncurses.a", pc.ProjectRoot+"/libncurses.a")
	return nil
}

func (a *ConfigureAdapter) Targets(pc *config.ProjectConfig) (targets []string, err error) {
	return nil, nil
}

func (a *ConfigureAdapter) CmakeConfig(pc *config.ProjectConfig) (string, error) {
	return "", nil
}
