package sdladapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type SdlAdapter struct {
}

// New creates a new builtin MakefileAdapter
func New() *SdlAdapter {
	return &SdlAdapter{}
}

var CurrentVersion, _ = manifest.VersionFromString("0.0.1")

func visit(path string, old string, new string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	newContents := strings.Replace(string(read), old, new, -1)
	err = ioutil.WriteFile(path, []byte(newContents), 0)
	if err != nil {
		return err
	}
	return nil
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


func createdDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func buildOnMacOS(pc *config.ProjectConfig) error {
	pathBuild := filepath.Join(pc.ProjectRoot, "build")
	err := createdDirectory(pathBuild)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = executeCli(filepath.Join(pathBuild, "../configure"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = executeCli("make")
	if err != nil {
		log.Fatal(err)
		return err
	}
	//err = executeCli( "make", "install")
	//if err != nil {
	//	log.Fatal(err)
	//	return err
	//}
	return nil
}

func buildOnLinux(pc *config.ProjectConfig) error {
	return nil
}

func (a *SdlAdapter) Build(pc *config.ProjectConfig) error {
	switch runtime.GOOS {
	case "darwin":
		err := buildOnMacOS(pc)
		if err != nil {
			return err
		}
		oldLocation := pc.ProjectRoot + "/build/.libs/libSDL2.dylib"
		err = os.Rename(oldLocation, pc.ProjectRoot+"/libSDL2.dylib")
		if err != nil {
			return err
		}
	case "linux":
		err := buildOnLinux(pc)
		if err != nil {
			return err
		}
	case "windows":
		return nil
	}
	return nil
}

func (a *SdlAdapter) CmakeConfig(_ *config.ProjectConfig) (string, error) {
	return "", nil
}

func (a *SdlAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}
