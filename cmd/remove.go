package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [dependencies]",
	Short: "Remove one or more dependency",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dependencies := args
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		err = ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: dependencies})
		if err != nil {
			return fmt.Errorf("failed to remove dependencies: %w", err)
		}
		return nil
	},
}
