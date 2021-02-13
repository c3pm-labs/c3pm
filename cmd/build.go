package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
)

type BuildCmd struct{}

func (b *BuildCmd) Run() error {
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return ctpm.AddDependenciesAndBuild(pc)
}
