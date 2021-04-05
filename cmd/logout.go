package cmd

import (
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use: "logout",
	Short: "Logout from the api",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ctpm.Logout()
	},
}
