package cmd

import "github.com/c3pm-labs/c3pm/ctpm"

//LogoutCmd defines the parameters of the logout command.
type LogoutCmd struct{}

//Run handles the behavior of the logout command.
func (l *LogoutCmd) Run() error {
	return ctpm.Logout()
}
