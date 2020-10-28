package cmd

import (
	"github.com/gabrielcolson/c3pm/cli/api"
	"github.com/gabrielcolson/c3pm/cli/cmd/input"
	"github.com/gabrielcolson/c3pm/cli/ctpm"
	"net/http"
)

type LoginCmd struct {
}

func (l *LoginCmd) Run() error {
	payload, err := input.Login()
	if err != nil {
		return err
	}
	client := api.API{Client: &http.Client{}}
	return ctpm.Login(client, payload.Login, payload.Password)
}
