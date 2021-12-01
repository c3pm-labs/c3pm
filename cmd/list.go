package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

type ListCmdFlags struct {
	ctpm.ListOptions
}

var listCmdFlags = ListCmdFlags{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all project dependencies",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		err = ctpm.List(pc, listCmdFlags.ListOptions)
		if err != nil {
			return fmt.Errorf("failed to add dependencies: %w", err)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listCmdFlags.Tree, "tree", ctpm.ListDefaultOptions.Tree, "List dependencies in indented tree form")
}
