package defaultadapter

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/internal/cmake"
	"path/filepath"
)

// DefaultAdapter is the builtin adapter used by default in c3pm
type DefaultAdapter struct {
}

// New creates a new builtin DefaultAdapter
func New() *DefaultAdapter {
	return &DefaultAdapter{}
}

var CurrentVersion, _ = manifest.VersionFromString("0.0.1")

func (a *DefaultAdapter) Build(pc *config.ProjectConfig) error {
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

	err = generateCMakeScripts(cmakeDirFromPc(pc), pc)
	if err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	err = cmake.GenerateBuildFiles(cmakeDirFromPc(pc), buildDirFromPc(pc), cmakeVariables)
	if err != nil {
		return fmt.Errorf("cmake build failed: %w", err)
	}

	err = cmake.Build(buildDirFromPc(pc))
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

	hasSources, err := hasFileMatchingRule(cfg.Sources, pc.ProjectRoot)
	return !hasSources, err
}

func hasFileMatchingRule(rules []string, projectRoot string) (bool, error) {
	for _, rule := range rules {
		files, err := doublestar.Glob(filepath.Join(projectRoot, rule))
		if err != nil {
			return false, err
		}
		if len(files) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (a *DefaultAdapter) Targets(_ *config.ProjectConfig) ([]string, error) {
	return nil, nil
}

func cmakeDirFromPc(pc *config.ProjectConfig) string {
	return filepath.Join(pc.LocalC3PMDirPath(), "cmake")
}

func buildDirFromPc(pc *config.ProjectConfig) string {
	return filepath.Join(pc.LocalC3PMDirPath(), "build")
}
