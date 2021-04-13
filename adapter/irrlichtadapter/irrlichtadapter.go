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

func executeXcodeCli(args ...string) error {
	cmd := exec.Command("xcodebuild", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start cmake: %w", err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("cmake process failed: %w", err)
	}
	return nil
}

func buildMacOS(pc *config.ProjectConfig) error {
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
	return executeXcodeCli("-project", path, "-target", "libIrrlicht.a", "SYSMROOT=build")
}

func (a *IrrlichtAdapter) Build(pc *config.ProjectConfig) error {
	if runtime.GOOS == "darwin" {
		err := buildMacOS(pc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *IrrlichtAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}
