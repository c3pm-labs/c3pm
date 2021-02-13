package cmakegen

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Dependency struct {
	Name                string
	Version             string
	Targets             []string
	ExportedDir         string
	ExportedIncludeDirs []string
}

type CMakeVars struct {
	ProjectName      string
	ProjectVersion   string
	Sources          string
	Includes         string
	IncludeDirs      string
	ExportedDir      string
	C3pmGlobalDir    string
	Dependencies     []Dependency
	PublicIncludeDir string
	LinuxConfig      *manifest.LinuxConfig
	LanguageStandard string
}

func dependenciesToCmake(dependencies map[string]string) ([]Dependency, error) {
	deps := make([]Dependency, len(dependencies))
	i := 0
	for n, v := range dependencies {
		m, err := manifest.Load(filepath.Join(config.LibCachePath(n, v), "c3pm.yml"))
		if err != nil {
			continue
		}
		deps[i] = Dependency{
			Name:                n,
			Version:             v,
			Targets:             m.Targets(),
			ExportedDir:         m.Files.ExportedDir,
			ExportedIncludeDirs: m.Files.ExportedIncludeDirs,
		}
		if m.CustomCmake != nil {
			deps[i].ExportedIncludeDirs = []string{"include"}
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
	dependencies, err := dependenciesToCmake(pc.Manifest.Dependencies)
	if err != nil {
		return CMakeVars{}, err
	}

	vars := CMakeVars{
		ProjectName:      pc.Manifest.Name,
		ProjectVersion:   pc.Manifest.Version.String(),
		Sources:          filesSliceToCMake(pc.Manifest.Files.Sources),
		Includes:         filesSliceToCMake(pc.Manifest.Files.Includes),
		IncludeDirs:      filesSliceToCMake(pc.Manifest.Files.IncludeDirs),
		ExportedDir:      filepath.ToSlash(filepath.Join(pc.ProjectRoot, pc.Manifest.Files.ExportedDir)),
		C3pmGlobalDir:    filepath.ToSlash(config.GlobalC3pmDirPath()),
		Dependencies:     dependencies,
		LinuxConfig:      pc.Manifest.LinuxConfig,
		LanguageStandard: pc.Manifest.Standard,
	}

	vars.Sources, err = globbingExprsToCMakeVar(pc.Manifest.Files.Sources, pc.ProjectRoot)
	if err != nil {
		return CMakeVars{}, fmt.Errorf("could not parse Sources: %w", err)
	}
	vars.Includes, err = globbingExprsToCMakeVar(pc.Manifest.Files.Includes, pc.ProjectRoot)
	if err != nil {
		return CMakeVars{}, fmt.Errorf("could not parse Includes: %w", err)
	}
	vars.IncludeDirs, err = globbingExprsToCMakeVar(pc.Manifest.Files.IncludeDirs, pc.ProjectRoot)
	if err != nil {
		return CMakeVars{}, fmt.Errorf("could not parse IncludeDirs: %w", err)
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

func Generate(pc *config.ProjectConfig) error {
	cmakeContent, err := fromProjectConfig(pc)
	if err != nil {
		return fmt.Errorf("failed to generate cmake scripts: %w", err)
	}
	err = os.MkdirAll(pc.CMakeDir(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create c3pm cmake directory: %w", err)
	}
	err = ioutil.WriteFile(filepath.Join(pc.CMakeDir(), "CMakeLists.txt"), []byte(cmakeContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create CMakeLists.txt: %w", err)
	}
	return nil
}
