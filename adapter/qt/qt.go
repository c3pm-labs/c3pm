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

// Working
func execConfigure(pc *config.ProjectConfig, execPath string, path string) error {
	cmd := exec.Command("/bin/sh", filepath.Join(execPath, "configure"))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
}

func initRepository(pc *config.ProjectConfig, path string) error {
	fmt.Println("PATH: " + filepath.Join(path, "init-repository"))
	cmd := exec.Command("perl", filepath.Join(path, "init-repository"))
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
	// var path = filepath.Join(pc.ProjectRoot, "qt5")
	var pathBuild = filepath.Join(pc.ProjectRoot, "qt6-build")
	// err := initRepository(pc, path)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// err := os.Mkdir(pathBuild, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// err1 := execConfigure(pc, path, pathBuild)
	// if err1 != nil {
	// 	return err1
	// }
	err2 := execCmake(pathBuild)
	if err2 != nil {
		return err2
	}
	return nil
	// return execCmakeInstall(pathBuild)
}

func execCmakeInstall(path string, args ...string) error {
	fmt.Println("PATH: " + filepath.Join(path, "init-repository"))

	cmd := exec.Command("cmake", "--install", filepath.Join(path, "."))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
	// cmd := exec.Command("cmake", args...)
	// cmd.Dir = path
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Start()
	// if err != nil {
	// 	return fmt.Errorf("failed to start cmake: %w", err)
	// }
	// if err = cmd.Wait(); err != nil {
	// 	return fmt.Errorf("cmake process failed: %w", err)
	// }
	// return nil
}

func execCmake(path string) error {
	cmd := exec.Command("cmake", "--build", filepath.Join(path, "."))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Start()
	// if err != nil {
	// 	return fmt.Errorf("failed to start cmake: %w", err)
	// }
	// if err = cmd.Wait(); err != nil {
	// 	return fmt.Errorf("cmake process failed: %w", err)
	// }
	// return nil
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

func (a *QtAdapter) CmakeConfig(_ *config.ProjectConfig) (string, error) {
	return "", nil
}

func (a *QtAdapter) Targets(pc *config.ProjectConfig) (targets []string, err error) {
	return nil, nil
}

func NewAdapter() *QtAdapter {
	return &QtAdapter{}
}
