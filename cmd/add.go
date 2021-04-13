package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

var addCmdFlags = ctpm.AddOptions{}

var addCmd = &cobra.Command{
	Use:   "add [dependencies...]",
	Short: "Add one or more new dependency",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		addCmdFlags.Dependencies = args
		err = ctpm.Add(pc, addCmdFlags)
		if err != nil {
			return fmt.Errorf("failed to add dependencies: %w", err)
		}
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addCmdFlags.Force, "force", "f", false, "Ignore cache.")
	addCmd.Flags().StringVarP(&addCmdFlags.RegistryURL, "registry-url", "r", "", "Select specific registry to use.")
}
