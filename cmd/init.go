package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	"path/filepath"
)

type InitCmd struct {
	NoTemplate  bool   `kong:"optional,name='no-template',help='Prevents the creation of CMake files'"`
	Path        string `kong:"optional,arg,name='path',help='Project path, default to working directory',default='.'"`
	Name        string `kong:"optional,group='ni',name='name',help='Give project name to skip interactive entry and enter non-interactive mode'"`
	Executable  bool   `kong:"optional,xor='type',group='ni',name='executable',help='Specify if project is an executable when using non-interactive mode', default='True'"`
	Library     bool   `kong:"optional,xor='type',group='ni',name='library',help='Specify if project is a library when using non-interactive mode'"`
	Description string `kong:"optional,group='ni',name='description',help='Project description when using non-interactive mode'"`
	Version     string `kong:"optional,group='ni',name='pversion',help='Project version when using non-interactive mode',default='1.0.0'"`
	License     string `kong:"optional,group='ni',name='license',help='Project license when using non-interactive mode',default='UNLICENSED'"`
}

func (i *InitCmd) Run() error {
	var man manifest.Manifest
	var err error
	if len(i.Name) == 0 {
		man, err = input.Init()
	} else {
		man, err = input.InitNonInteractive(input.InitValues{
			Name:        i.Name,
			Executable:  i.Executable,
			Library:     i.Library,
			Description: i.Description,
			Version:     i.Version,
			License:     i.License,
		})
	}
	if err != nil {
		return fmt.Errorf("failed to init project config: %w", err)
	}
	projectRoot, err := filepath.Abs(i.Path)
	if err != nil {
		return err
	}
	pc := &config.ProjectConfig{Manifest: man, ProjectRoot: projectRoot}
	err = ctpm.Init(pc, ctpm.InitOptions{NoTemplate: i.NoTemplate})
	if err != nil {
		return fmt.Errorf("failed to init project: %w", err)
	}
	return nil
}
