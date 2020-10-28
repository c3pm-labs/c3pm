package cmd

import (
	"fmt"
	"github.com/gabrielcolson/c3pm/cli/config"
	"github.com/gabrielcolson/c3pm/cli/ctpm"
)

type BuildCmd struct{}

func (b *BuildCmd) Run() error {
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return ctpm.Build(pc)
}
