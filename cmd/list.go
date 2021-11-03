package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all project dependencies",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		err = ctpm.List(pc)
		if err != nil {
			return fmt.Errorf("failed to add dependencies: %w", err)
		}
		return nil
	},
}
