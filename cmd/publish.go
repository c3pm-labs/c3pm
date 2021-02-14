package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
	"net/http"
)

//PublishCmd defines the parameters of the publish command.
type PublishCmd struct {
}

//Run handles the behavior of the publish command.
func (p *PublishCmd) Run() error {
	token, err := config.TokenStrict()
	if err != nil {
		return fmt.Errorf("not logged in: %w", err)
	}
	client := api.New(&http.Client{}, token)
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return ctpm.Publish(pc, client, ctpm.PublishOptions{
		Exclude: pc.Manifest.Exclude,
		Include: pc.Manifest.Include,
	})
}
