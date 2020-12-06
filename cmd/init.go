package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
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
