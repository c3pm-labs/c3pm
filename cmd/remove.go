package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
)

//RemoveCmd defines the parameters of the remove command.
type RemoveCmd struct {
	Dependencies []string `kong:"arg,help='List of dependencies to remove.'"`
}

//Run handles the behavior of the remove command.
func (a *RemoveCmd) Run() error {
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	err = ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: a.Dependencies})
	if err != nil {
		return fmt.Errorf("failed to remove dependencies: %w", err)
	}
	return nil
}
