package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
	"path/filepath"
)

type InitCmdFlags struct {
	ctpm.InitOptions
	input.InitValues
}

var initCmdFlags = InitCmdFlags{}

var initCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Init a c3pm project",
	Long: "Init a c3pm project\n\n" +
		"Project path defaults to working directory",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("requires one or no arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if len(args) == 1 {
			path = args[0]
		} else {
			path = "."
		}
    
		var man manifest.Manifest
		var err error
		if len(initCmdFlags.Name) == 0 {
			man, err = input.Init()
		} else {
			man, err = input.InitNonInteractive(initCmdFlags.InitValues)
		}

    
		if err != nil {
			return fmt.Errorf("failed to init project config: %w", err)
		}
		projectRoot, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		pc := &config.ProjectConfig{Manifest: man, ProjectRoot: projectRoot}
		err = ctpm.Init(pc, initCmdFlags.InitOptions)
		if err != nil {
			return fmt.Errorf("failed to init project: %w", err)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initCmdFlags.NoTemplate, "no-template", ctpm.InitDefaultOptions.NoTemplate, "Prevents the creation of CMake files")

	// Non Interactive Mode
	initCmd.Flags().StringVar(&initCmdFlags.InitValues.Name, "name", "", "Give project name to skip interactive entry and enter non-interactive mode")
	initCmd.Flags().StringVar(&initCmdFlags.InitValues.Type, "type", "", "Project's type when using non-interactive mode")
	initCmd.Flags().StringVar(&initCmdFlags.InitValues.Description, "desc", "", "Project description when using non-interactive mode")
	initCmd.Flags().StringVar(&initCmdFlags.InitValues.Version, "version", "1.0.0", "Project version when using non-interactive mode")
	initCmd.Flags().StringVar(&initCmdFlags.InitValues.License, "license", "UNLICENSED", "Project license when using non-interactive mode")
}
