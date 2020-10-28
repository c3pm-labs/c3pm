package cmd

import "github.com/gabrielcolson/c3pm/cli/ctpm"

type LogoutCmd struct{}

func (l *LogoutCmd) Run() error {
	return ctpm.Logout()
}
