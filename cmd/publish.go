package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
	"net/http"
)

var publishCmdFlags = ctpm.PublishOptions{}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a c3pm project",
	Long: "Publish a c3pm project\n\n" +
		"[Ignore files]",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := config.TokenStrict()
		if err != nil {
			return fmt.Errorf("not logged in: %w", err)
		}
		client := api.New(&http.Client{}, token)
		pc, err := config.Load(".")
		if err != nil {
			return fmt.Errorf("failed to read c3pm.yml: %w", err)
		}
		return ctpm.Publish(pc, client, publishCmdFlags)
	},
}

func init() {
	publishCmd.Flags().StringArrayVarP(&publishCmdFlags.Ignore, "ignore", "i", ctpm.PublishDefaultOptions.Ignore, "Ignore files")
}
