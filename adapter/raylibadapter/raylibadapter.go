package raylibadapter

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
)

type RaylibAdapter struct {
}

// New creates a new builtin MakefileAdapter
func New() *RaylibAdapter {
	return &RaylibAdapter{}
}

var CurrentVersion, _ = manifest.VersionFromString("0.0.1")

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

func buildOnMacOS(pc *config.ProjectConfig) error {
	return executeCli("make", "-C", pc.ProjectRoot+"/src")
}

func buildOnLinux(pc *config.ProjectConfig) error {
	executeCli("mkdir", pc.ProjectRoot+"/build")
	executeCli("cd " + pc.ProjectRoot+"/build" + " && cmake -DBUILD_SHARED_LIBS=ON ..")
	return executeCli("make", pc.ProjectRoot+"/build")
}

func (a *RaylibAdapter) Build(pc *config.ProjectConfig) error {
	switch runtime.GOOS {
	case "darwin":
		err := buildOnMacOS(pc)
		if err != nil {
			return err
		}
		oldLocation := pc.ProjectRoot + "/src/libraylib.a"
		err = os.Rename(oldLocation, pc.ProjectRoot+"/libraylib.a")
		if err != nil {
			return err
		}
	case "linux":
		err := buildOnLinux(pc)
		if err != nil {
			return err
		}
	// case "windows":
	// 	return nil
	}
	return nil
}

func (a *RaylibAdapter) CmakeConfig(pc *config.ProjectConfig) (string, error) {
	return CmakeConfig, nil
}

func (a *RaylibAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}
