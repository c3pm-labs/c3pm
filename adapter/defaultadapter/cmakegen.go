// Package cmakegen handles the templating and generation of CMake configuration files.
package defaultadapter

import (
	"errors"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/adapter/irrlichtadapter"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//dependency is holds metadata about a dependency of a project.
type dependency struct {
	// Name is the package name of the dependency
	Name string
	// Version is the version of the dependency to depend on
	Version string
	// Targets is the list of the libraries contained by the dependencies.
	// In most cases this will only contain one entry, but there are cases of packages containing several libraries, for
	// separation of concerns reasons.
	Targets []string
	// IncludeDirs is the list of the directories in which header files for the library can be found.
	IncludeDirs []string
}

//cmakeVars is the structure passed to the templates used for generating CMake config files.
type cmakeVars struct {
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
	//Dependencies is a list of all the data for each dependency of the project
	Dependencies []dependency
	//TODO: Unused
	PublicIncludeDir string
	//LinuxConfig holds linux-specific configuration information
	LinuxConfig *LinuxConfig
	//LanguageStandard is the C++ language standard version to use.
	LanguageStandard string
	//DependenciesConfig is a string containing all the cmake command needed by dependencies
	DependenciesConfig string
}

// FIXME find a way to use adapter file
type Adapter interface {
	// Build builds the targets
	Build(pc *config.ProjectConfig) error
	// Targets return the paths of the targets built by the Build function
	Targets(pc *config.ProjectConfig) (targets []string, err error)
	CmakeConfig(pc *config.ProjectConfig) (string, error)
}

// FIXME find a way to use adapter file
func fromPC(adp *manifest.AdapterConfig) (Adapter, error) {

	switch {
	case adp.Name == "c3pm" && adp.Version.String() == "0.0.1":
		return New(), nil
	case adp.Name == "irrlicht" && adp.Version.String() == "0.0.1":
		return irrlichtadapter.New(), nil
	default:
		return nil, errors.New("only default adapter is supported")
	}
}

func dependenciesToCMake(pc *config.ProjectConfig) ([]dependency, string, error) {
	deps := make([]dependency, len(pc.Manifest.Dependencies))
	var depsConfig = ""
	i := 0
	for n, v := range pc.Manifest.Dependencies {
		m, err := manifest.Load(filepath.Join(config.LibCachePath(n, v), "c3pm.yml"))
		if err != nil {
			return nil, "", err
		}
		deps[i] = dependency{
			Name:        n,
			Version:     v,
			Targets:     m.Targets(),
			IncludeDirs: m.Publish.IncludeDirs,
		}
		adp, err := fromPC(m.Build.Adapter)
		if err != nil {
			return nil, "", err
		}
		dependencyConfig, err := adp.CmakeConfig(pc)
		if err != nil {
			return nil, "", err
		}
		depsConfig = depsConfig + dependencyConfig
	}
	return deps, depsConfig, nil
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
	return strings.Join(files, " "), nil
}

func pathListToCmakeVar(paths []string, projectRoot string) string {
	res := ""
	for _, path := range paths {
		res += " "
		res += filepath.Join(projectRoot, path)
	}
	return res
}

func varsFromProjectConfig(pc *config.ProjectConfig) (cmakeVars, error) {
	dependencies, dependenciesConfig, err := dependenciesToCMake(pc)
	if err != nil {
		return cmakeVars{}, err
	}

	adapterCfg, err := Parse(pc.Manifest.Build.Config)
	if err != nil {
		return cmakeVars{}, err
	}

	sources, err := globbingExprsToCMakeVar(adapterCfg.Sources, pc.ProjectRoot)
	if err != nil {
		return cmakeVars{}, fmt.Errorf("could not parse Sources: %w", err)
	}
	headers, err := globbingExprsToCMakeVar(adapterCfg.Headers, pc.ProjectRoot)
	if err != nil {
		return cmakeVars{}, fmt.Errorf("could not parse Includes: %w", err)
	}

	vars := cmakeVars{
		ProjectName:        pc.Manifest.Name,
		ProjectVersion:     pc.Manifest.Version.String(),
		Sources:            sources,
		Headers:            headers,
		IncludeDirs:        pathListToCmakeVar(adapterCfg.IncludeDirs, pc.ProjectRoot),
		C3PMGlobalDir:      filepath.ToSlash(config.GlobalC3PMDirPath()),
		Dependencies:       dependencies,
		LinuxConfig:        adapterCfg.LinuxConfig,
		LanguageStandard:   pc.Manifest.Standard,
		DependenciesConfig: dependenciesConfig,
	}

	return vars, nil
}

func fromProjectConfig(pc *config.ProjectConfig) (string, error) {
	var cmake string
	var vars cmakeVars

	vars, err := varsFromProjectConfig(pc)
	if err != nil {
		return "", fmt.Errorf("failed to generate cmake variables: %w", err)
	}
	switch pc.Manifest.Type {
	case manifest.Executable:
		cmake, err = (func() (string, error) { return executable(vars) })()
	case manifest.Library:
		cmake, err = (func() (string, error) { return library(vars) })()
	}
	if err != nil {
		return "", fmt.Errorf("failed to generate cmake: %w", err)
	}
	return cmake, nil
}

//generateCMakeScripts takes a config.ProjectConfig and creates CMake configuration files based on the project config.
func generateCMakeScripts(targetDir string, pc *config.ProjectConfig) error {
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
