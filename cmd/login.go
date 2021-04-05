package cmd

import (
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/ctpm"
	"github.com/spf13/cobra"
	"net/http"
)

var loginCmd = &cobra.Command{
	Use: "login",
	Short: "Login to the api",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		payload, err := input.Login()
		if err != nil {
			return err
		}
		client := api.New(&http.Client{}, "")
		return ctpm.Login(client, payload.Login, payload.Password)
	},
}
