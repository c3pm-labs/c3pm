package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmake"
	"github.com/c3pm-labs/c3pm/cmakegen"
	"github.com/c3pm-labs/c3pm/config"
	"path/filepath"
)

func Build(pc *config.ProjectConfig) error {
	cmakeVariables := map[string]string{
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY": pc.ProjectRoot,
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY": pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY": pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_INSTALL_PREFIX":           filepath.ToSlash(filepath.Join(config.GlobalC3pmDirPath(), "cache", pc.Manifest.Name, pc.Manifest.Version.String())),
		"CMAKE_BUILD_TYPE":               "Release",
		// Useful for Windows build
		//"MSVC_TOOLSET_VERSION":           "141",
		//"MSVC_VERSION":                   "1916",
	}

	if pc.UseCustomCmake() {
		for key, value := range pc.Manifest.CustomCmake.Variables {
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

	err = cmake.Build(pc.BuildDir())
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	return nil
}
