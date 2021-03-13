package builtin

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/adapter/builtin/cmake"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"path/filepath"
)

// Adapter is the builtin adapter used by default in c3pm
type Adapter struct {
}

// checks if Adapter implements the adapter.Adapter interface
var _ adapter.Adapter = (*Adapter)(nil)

// New creates a new builtin Adapter
func New() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Build(pc *config.ProjectConfig) error {
	cmakeVariables := map[string]string{
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_LIBRARY_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY":         pc.ProjectRoot,
		"CMAKE_RUNTIME_OUTPUT_DIRECTORY_RELEASE": pc.ProjectRoot,
		"CMAKE_INSTALL_PREFIX":                   filepath.ToSlash(filepath.Join(config.GlobalC3PMDirPath(), "cache", pc.Manifest.Name, pc.Manifest.Version.String())),
		"CMAKE_BUILD_TYPE":                       "Release",
		// Useful for Windows build
		//"MSVC_TOOLSET_VERSION":           "141",
		//"MSVC_VERSION":                   "1916",
	}

	headerOnly, err := isHeaderOnly(pc)
	if err != nil {
		return err
	}

	// don't build if the lib is header only
	if headerOnly && pc.Manifest.Type == manifest.Library {
		// TODO: generate cmake files so we can have IDE integration
		return nil
	}

	err = GenerateScripts(CMakeDirFromPc(pc), pc)
	if err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	err = cmake.GenerateBuildFiles(CMakeDirFromPc(pc), BuildDirFromPc(pc), cmakeVariables)
	if err != nil {
		return fmt.Errorf("cmake build failed: %w", err)
	}

	err = cmake.Build(BuildDirFromPc(pc))
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	return nil
}

func isHeaderOnly(pc *config.ProjectConfig) (bool, error) {
	cfg, err := Parse(pc.Manifest.Build.Config)
	if err != nil {
		return false, err
	}

	hasSources, err := hasFileMatchingRule(cfg.Sources)
	return !hasSources, err
}

func hasFileMatchingRule(rules []string) (bool, error) {
	for _, rule := range rules {
		files, err := doublestar.Glob(rule)
		if err != nil {
			return false, err
		}
		if len(files) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (a *Adapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}

func CMakeDirFromPc(pc *config.ProjectConfig) string {
	return filepath.Join(pc.LocalC3PMDirPath(), "cmake")
}

func BuildDirFromPc(pc *config.ProjectConfig) string {
	return filepath.Join(pc.LocalC3PMDirPath(), "build")
}
