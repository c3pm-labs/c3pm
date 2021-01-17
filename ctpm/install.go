package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmake"
	"github.com/c3pm-labs/c3pm/cmakegen"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"path/filepath"
)

// TODO: unused
func Install(pc *config.ProjectConfig) error {
	libDir := filepath.Join(config.GlobalC3PMDirPath(), "cache", pc.Manifest.Name, pc.Manifest.Version.String())

	cmakeVariables := map[string]string{
		"CMAKE_INSTALL_BINDIR": filepath.Join(config.GlobalC3PMDirPath(), "bin"),
		"CMAKE_INSTALL_PREFIX": libDir,
	}

	if pc.UseCustomCMake() {
		for key, value := range pc.Manifest.CustomCMake.Variables {
			cmakeVariables[key] = value
		}
	} else {
		err := cmakegen.Generate(pc)
		if err != nil {
			return fmt.Errorf("error generating config files: %w", err)
		}
	}

	err := cmake.GenerateBuildFiles(pc.CMakeDir(), pc.BuildDir(), cmakeVariables)
	if err != nil {
		return fmt.Errorf("cmake build failed: %w", err)
	}

	err = cmake.Install(pc.BuildDir())
	if err != nil {
		return fmt.Errorf("install failed: %w", err)
	}

	if pc.Manifest.Type == manifest.Library {
		return pc.Manifest.Save(filepath.Join(libDir, "c3pm.yml"))
	}

	return nil
}
