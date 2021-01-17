// Package config handles the interactions with C3PM's various configuration files.
// It handles interaction with both the c3pm.yml file (see package manifest), and
// the storage of authentication tokens in the global C3PM directory as found by GlobalC3pmDirPath.
package config

import (
	"github.com/c3pm-labs/c3pm/config/manifest"
	"os"
	"path"
	"path/filepath"
)

//ProjectConfig represents the configuration of a C3PM project.
type ProjectConfig struct {
	//Manifest is the representation of the contents of the c3pm.yml file.
	Manifest manifest.Manifest
	//ProjectRoot stores the absolute path to the C3PM project.
	ProjectRoot string
}

//Load takes the path of a project and creates the ProjectConfig object that represents the configuration of this project.
func Load(projectPath string) (*ProjectConfig, error) {
	absolutePath, err := filepath.Abs(projectPath)
	if err != nil {
		return nil, err
	}

	m, err := manifest.Load(filepath.Join(absolutePath, "c3pm.yml"))
	if err != nil {
		return nil, err
	}

	return &ProjectConfig{
		Manifest:    m,
		ProjectRoot: absolutePath,
	}, nil
}

//Save writes the current configuration and writes it to the project directory.
func (pc *ProjectConfig) Save() error {
	return pc.Manifest.Save(filepath.Join(pc.ProjectRoot, "c3pm.yml"))
}

//BuildDir returns the path to the build directory used for CMake build files.
func (pc *ProjectConfig) BuildDir() string {
	return filepath.Join(pc.ProjectRoot, ".c3pm", "build")
}

//CMakeDir returns the path to the CMake files to use for the project.
func (pc *ProjectConfig) CMakeDir() string {
	if pc.UseCustomCmake() {
		return filepath.Join(pc.ProjectRoot, pc.Manifest.CustomCmake.Path)
	}
	return filepath.Join(pc.ProjectRoot, ".c3pm", "cmake")
}

//GlobalC3pmDirPath finds the path to the global C3PM directory.
func GlobalC3pmDirPath() string {
	if dir := os.Getenv("C3PM_USER_DIR"); dir != "" {
		return dir
	}
	homeDir := os.Getenv("HOME")
	return path.Join(homeDir, ".c3pm")
}

//UseCustomCmake checks if a custom CMake configuration is to be used for the project.
func (pc *ProjectConfig) UseCustomCmake() bool {
	return pc.Manifest.CustomCmake != nil
}

//LibCachePath returns the path to the global C3PM cache
func LibCachePath(name, version string) string {
	return filepath.Join(GlobalC3pmDirPath(), "cache", name, version)
}
