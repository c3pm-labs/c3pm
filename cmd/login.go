package cmd

import (
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/cmd/input"
	"github.com/c3pm-labs/c3pm/ctpm"
	"net/http"
)

type LoginCmd struct {
}

func (l *LoginCmd) Run() error {
	payload, err := input.Login()
	if err != nil {
		return err
	}
	client := api.New(&http.Client{}, "")
	return ctpm.Login(client, payload.Login, payload.Password)
}
