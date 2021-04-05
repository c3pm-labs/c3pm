package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use: "build",
	Short: "Build a c3pm project",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		return ctpm.AddDependenciesAndBuild(pc)
	},
}
