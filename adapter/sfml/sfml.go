package sfml

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/internal/cmake"
	"path/filepath"
)

type Adapter struct{}

func (a *Adapter) Build(pc *config.ProjectConfig) error {
	cmakeVariables := map[string]string{
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_BUILD_TYPE":                       "Release",
		"BUILD_SHARED_LIB":                       "OFF",
	}

	sourceDir := pc.ProjectRoot
	buildDir := filepath.Join(pc.ProjectRoot, "build")

	err := cmake.GenerateBuildFiles(sourceDir, buildDir, cmakeVariables)
	if err != nil {
		return fmt.Errorf("cmake build failed: %w", err)
	}

	err = cmake.Build(buildDir)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	return nil
}

func (a *Adapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return []string{
		"sfml-system",
		"sfml-window",
		"sfml-graphics",
		"sfml-network",
		"sfml-audio",
	}, nil
}

func (a *Adapter) CmakeConfig(_ *config.ProjectConfig) (string, error) {
	return "", nil
}

func New() *Adapter {
	return &Adapter{}
}
