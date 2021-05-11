package irrlichtadapter

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type IrrlichtAdapter struct {
}

// New creates a new builtin MakefileAdapter
func New() *IrrlichtAdapter {
	return &IrrlichtAdapter{}
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

func buildOnMacOS(pc *config.ProjectConfig) error {
	err := visit(pc.ProjectRoot+"/src/Irrlicht/MacOSX/CIrrDeviceMacOSX.mm",
		"[NSApp setDelegate:(id<NSFileManagerDelegate>)",
		"[NSApp setDelegate:(id<NSApplicationDelegate>)",
	)
	if err != nil {
		return err
	}
	err = visit(pc.ProjectRoot+"/src/Irrlicht/libpng/pngpriv.h",
		"#  error ZLIB_VERNUM != PNG_ZLIB_VERNUM \\",
		"#  warning ZLIB_VERNUM != PNG_ZLIB_VERNUM \\",
	)
	if err != nil {
		return err
	}
	var path = pc.ProjectRoot + "/src/Irrlicht/MacOSX/MacOSX.xcodeproj"
	return executeCli("xcodebuild", "-project", path, "-target", "libIrrlicht.a", "SYSMROOT=build")
}

func buildOnLinux(pc *config.ProjectConfig) error {
	return executeCli("make", "-C", pc.ProjectRoot+"/src/Irrlicht")
}

func (a *IrrlichtAdapter) Build(pc *config.ProjectConfig) error {
	switch runtime.GOOS {
	case "darwin":
		err := buildOnMacOS(pc)
		if err != nil {
			return err
		}
		oldLocation := pc.ProjectRoot + "/src/Irrlicht/MacOSX/build/Release/libIrrlicht.a"
		err = os.Rename(oldLocation, pc.ProjectRoot+"/libIrrlicht.a")
		if err != nil {
			return err
		}
	case "linux":
		err := buildOnLinux(pc)
		if err != nil {
			return err
		}
		oldLocation := pc.ProjectRoot + "/lib/Linux/libIrrlicht.a"
		err = os.Rename(oldLocation, pc.ProjectRoot+"/libIrrlicht.a")
		if err != nil {
			return err
		}
	case "windows":
		return nil
	}
	return nil
}

func (a *IrrlichtAdapter) CmakeConfig(pc *config.ProjectConfig) (string, error) {
	return CmakeConfig, nil
}

func (a *IrrlichtAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return []string{"Irrlicht"}, nil
}
