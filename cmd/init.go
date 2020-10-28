package cmd

import (
	"fmt"
	"github.com/gabrielcolson/c3pm/cli/cmd/input"
	"github.com/gabrielcolson/c3pm/cli/config"
	"github.com/gabrielcolson/c3pm/cli/ctpm"
	"path/filepath"
)

type InitCmd struct {
	NoTemplate bool   `kong:"optional,name='no-template',help='Prevents the creation of CMake files'"`
	Path       string `kong:"optional,arg,name='path',help='Project path, default to working directory',default='.'"`
}

func (i *InitCmd) Run() error {
	manifest, err := input.Init()
	if err != nil {
		return fmt.Errorf("failed to init project config: %w", err)
	}
	projectRoot, err := filepath.Abs(i.Path)
	if err != nil {
		return err
	}
	pc := &config.ProjectConfig{Manifest: manifest, ProjectRoot: projectRoot}
	err = ctpm.Init(pc, ctpm.InitOptions{NoTemplate: i.NoTemplate})
	if err != nil {
		return fmt.Errorf("failed to init project: %w", err)
	}
	return nil
}
