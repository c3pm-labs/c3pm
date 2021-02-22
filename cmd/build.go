package cmd

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/ctpm"
)

//BuildCmd defines the parameters of the build command.
type BuildCmd struct{}

//Run handles the behavior of the build command.
func (b *BuildCmd) Run() error {
	pc, err := config.Load(".")
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return ctpm.AddDependenciesAndBuild(pc)
}
