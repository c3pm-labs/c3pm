package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
)

type AddCmd struct {
	Force       bool   `kong:"optional,name='force',help='Ignore cache.'"`
	RegistryURL string `kong:"optional,name='registry-url',help='Select specific registry to use.'"`

	Dependencies []string `kong:"arg,help='List of dependencies to add.'"`
}

func (a *AddCmd) Run() error {
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	err = ctpm.Add(pc, ctpm.AddOptions{Force: a.Force, RegistryURL: a.RegistryURL, Dependencies: a.Dependencies})
	if err != nil {
		return fmt.Errorf("Failed to add dependencies: %w", err)
	}
	return nil
}
