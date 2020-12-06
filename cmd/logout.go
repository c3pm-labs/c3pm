package cmd

import "github.com/c3pm-labs/c3pm/ctpm"

type LogoutCmd struct{}

func (l *LogoutCmd) Run() error {
	return ctpm.Logout()
}
