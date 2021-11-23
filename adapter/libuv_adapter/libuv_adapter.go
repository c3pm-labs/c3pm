package libuv_adapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"os"
	"os/exec"
	"runtime"
)

type UvAdapter struct {
}

// New creates a new builtin MakefileAdapter
func New() *UvAdapter {
	return &UvAdapter{}
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

func (a *UvAdapter) Build(pc *config.ProjectConfig) error {
	switch runtime.GOOS {
	case "windows":
		fmt.Println("not supported")
		break
	default:
		err := executeCli("cmake", "-S", pc.ProjectRoot, "-B", pc.ProjectRoot)
		if err != nil {
			return err
		}
		err = executeCli("make", "-C", pc.ProjectRoot)
		if err != nil {
			return err
		}
	}
	return nil
}


func (a *UvAdapter) CmakeConfig(pc *config.ProjectConfig) (string, error) {
	return "", nil
}

func (a *UvAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}

