package config

import (
	"github.com/gabrielcolson/c3pm/cli/config/manifest"
	"os"
	"path"
	"path/filepath"
)

type ProjectConfig struct {
	Manifest    manifest.Manifest
	ProjectRoot string
}

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

func (pc *ProjectConfig) Save() error {
	return pc.Manifest.Save(filepath.Join(pc.ProjectRoot, "c3pm.yml"))
}

func (pc *ProjectConfig) BuildDir() string {
	return filepath.Join(pc.ProjectRoot, ".c3pm", "build")
}

func (pc *ProjectConfig) CMakeDir() string {
	if pc.UseCustomCmake() {
		return filepath.Join(pc.ProjectRoot, pc.Manifest.CustomCmake.Path)
	}
	return filepath.Join(pc.ProjectRoot, ".c3pm", "cmake")
}

func GlobalC3pmDirPath() string {
	if dir := os.Getenv("C3PM_USER_DIR"); dir != "" {
		return dir
	}
	homeDir := os.Getenv("HOME")
	return path.Join(homeDir, ".c3pm")
}

func (pc *ProjectConfig) UseCustomCmake() bool {
	return pc.Manifest.CustomCmake != nil
}

func LibCachePath(name, version string) string {
	return filepath.Join(GlobalC3pmDirPath(), "cache", name, version)
}
