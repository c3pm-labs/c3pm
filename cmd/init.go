package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
	"path/filepath"
)

var initCmdFlags = ctpm.InitOptions{}

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
		manifest, err := input.Init()
		if err != nil {
			return fmt.Errorf("failed to init project config: %w", err)
		}
		projectRoot, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		pc := &config.ProjectConfig{Manifest: manifest, ProjectRoot: projectRoot}
		err = ctpm.Init(pc, initCmdFlags)
		if err != nil {
			return fmt.Errorf("failed to init project: %w", err)
		}
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initCmdFlags.NoTemplate, "no-template", ctpm.InitDefaultOptions.NoTemplate, "Prevents the creation of CMake files")
}
