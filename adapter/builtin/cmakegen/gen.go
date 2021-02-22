// Package cmakegen handles the templating and generation of CMake configuration files.
package cmakegen

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	builtinAdapterConfig "github.com/c3pm-labs/c3pm/adapter/builtin/config"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//Dependency is holds metadata about a dependency of a project.
type Dependency struct {
	// Name is the package name of the dependency
	Name string
	// Version is the version of the dependency to depend on
	Version string
	// Targets is the list of the libraries contained by the dependencies.
	// In most cases this will only contain one entry, but there are cases of packages containing several libraries, for
	// separation of concerns reasons.
	Targets []string
	//TODO: Unused
	ExportedDir string
	// ExportedIncludeDirs is the list of the directories in which header files for the library can be found.
	ExportedIncludeDirs []string
}

//CMakeVars is the structure passed to the templates used for generating CMake config files.
type CMakeVars struct {
	//ProjectName is the name of the current project
	ProjectName string
	//ProjectVersion is the current version of the project
	ProjectVersion string
	//Sources is a string containing the list of all of the project's sources, space-separated.
	Sources string
	//Headers is string containing the list of all of the project's header files, space-separated.
	Headers string
	//IncludeDirs is a string containing the list of all of the project's additional header directories, space-separated.
	IncludeDirs string
	//ExportedDir is the path to the directory containing export headers for the project.
	ExportedDir string
	//C3PMGlobalDir is the path to the current $HOME user directory.
	C3PMGlobalDir string
	//Dependencies is a list of all the data for each Dependency of the project
	Dependencies []Dependency
	//TODO: Unused
	PublicIncludeDir string
	//LinuxConfig holds linux-specific configuration information
	LinuxConfig *builtinAdapterConfig.LinuxConfig
	//LanguageStandard is the C++ language standard version to use.
	LanguageStandard string
}

func dependenciesToCMake(dependencies map[string]string) ([]Dependency, error) {
	deps := make([]Dependency, len(dependencies))
	i := 0
	for n, v := range dependencies {
		m, err := manifest.Load(filepath.Join(config.LibCachePath(n, v), "c3pm.yml"))
		if err != nil {
			return nil, err
		}
		deps[i] = Dependency{
			Name:    n,
			Version: v,
			Targets: m.Targets(),
			//ExportedDir:         m.Files.ExportedDir,
			//ExportedIncludeDirs: m.Files.ExportedIncludeDirs,
		}
	}
	return deps, nil
}

func filesSliceToCMake(files []string) string {
	fileString := ""
	for _, file := range files {
		fileString += " " + filepath.ToSlash(file)
	}
	return fileString
}

func globbingExprToFiles(globStr string) ([]string, error) {
	return doublestar.Glob(globStr)
}

func filterInternalSources(files []string, projectRoot string) []string {
	var newFiles []string
	for _, file := range files {
		if !strings.HasPrefix(file, filepath.Join(projectRoot, ".c3pm")) {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}

func globbingExprsToCMakeVar(globs []string, projectRoot string) (string, error) {
	var files []string
	for _, glob := range globs {
		globFiles, err := globbingExprToFiles(filepath.Join(projectRoot, glob))
		if err != nil {
			return "", fmt.Errorf("could not get files from globbing expression: %w", err)
		}
		files = append(files, globFiles...)
	}
	files = filterInternalSources(files, projectRoot)
	return filesSliceToCMake(files), nil
}

func varsFromProjectConfig(pc *config.ProjectConfig) (CMakeVars, error) {
	dependencies, err := dependenciesToCMake(pc.Manifest.Dependencies)
	if err != nil {
		return CMakeVars{}, err
	}

	adapterCfg, err := builtinAdapterConfig.Parse(pc.Manifest.Build.Config)
	if err != nil {
		return CMakeVars{}, err
	}

	sources, err := globbingExprsToCMakeVar(adapterCfg.Sources, pc.ProjectRoot)
	if err != nil {
		return CMakeVars{}, fmt.Errorf("could not parse Sources: %w", err)
	}
	headers, err := globbingExprsToCMakeVar(adapterCfg.Headers, pc.ProjectRoot)
	if err != nil {
		return CMakeVars{}, fmt.Errorf("could not parse Includes: %w", err)
	}
	//includeDirs, err := globbingExprsToCMakeVar(pc.Manifest.Files.IncludeDirs, pc.ProjectRoot)
	//if err != nil {
	//	return CMakeVars{}, fmt.Errorf("could not parse IncludeDirs: %w", err)
	//}

	vars := CMakeVars{
		ProjectName:    pc.Manifest.Name,
		ProjectVersion: pc.Manifest.Version.String(),
		Sources:        sources,
		Headers:        headers,
		//IncludeDirs:      includeDirs,
		//ExportedDir:      filepath.ToSlash(filepath.Join(pc.ProjectRoot, pc.Manifest.Files.ExportedDir)),
		C3PMGlobalDir: filepath.ToSlash(config.GlobalC3PMDirPath()),
		Dependencies:  dependencies,
		LinuxConfig:   adapterCfg.LinuxConfig,
		//LanguageStandard: pc.Manifest.Standard,
	}

	return vars, nil
}

func fromProjectConfig(pc *config.ProjectConfig) (string, error) {
	var cmake string
	var cmakeVars CMakeVars

	cmakeVars, err := varsFromProjectConfig(pc)
	if err != nil {
		return "", fmt.Errorf("failed to generate cmake variables: %w", err)
	}
	switch pc.Manifest.Type {
	case manifest.Executable:
		cmake, err = executable(cmakeVars)
	case manifest.Library:
		cmake, err = library(cmakeVars)
	}
	if err != nil {
		return "", fmt.Errorf("failed to generate cmake: %w", err)
	}
	return cmake, nil
}

//GenerateScripts takes a config.ProjectConfig and creates CMake configuration files based on the project config.
func GenerateScripts(targetDir string, pc *config.ProjectConfig) error {
	cmakeContent, err := fromProjectConfig(pc)
	if err != nil {
		return fmt.Errorf("failed to generate cmake scripts: %w", err)
	}
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create c3pm cmake directory: %w", err)
	}
	err = ioutil.WriteFile(filepath.Join(targetDir, "CMakeLists.txt"), []byte(cmakeContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt: %w", err)
	}
	return nil
}
