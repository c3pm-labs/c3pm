package qt

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/c3pm-labs/c3pm/config"
)

type QtAdapter struct{}

func New() *QtAdapter {
	return &QtAdapter{}
}

func execFile(pc *config.ProjectConfig, path string) error {
	fmt.Println("PATH: " + filepath.Join(path, "configure"))
	cmd := exec.Command("/bin/sh", filepath.Join(path, "configure"))
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
}

func executeMakeCli(path string, args ...string) error {
	cmd := exec.Command("make", args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start make: %w", err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("make process failed: %w", err)
	}
	return nil
}

func buildOnLinux(pc *config.ProjectConfig) error {
	var path = filepath.Join(pc.ProjectRoot, "qtbase")
	err := execFile(pc, path)
	if err != nil {
		return err
	}
	return executeMakeCli(path, "-C", path)
}

func (a *QtAdapter) Build(pc *config.ProjectConfig) error {
	fmt.Println(runtime.GOOS)
	switch runtime.GOOS {
	case "darwin":
		return nil
	case "linux":
		err := buildOnLinux(pc)
		if err != nil {
			return err
		}
		return err
	case "windows":
		return nil

	}
	return nil
}

func (a *QtAdapter) Targets(pc *config.ProjectConfig) (targets []string, err error) {
	return nil, nil
}

func NewAdapter() *QtAdapter {
	return &QtAdapter{}
}
