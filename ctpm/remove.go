package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
)

type RemoveOptions struct {
	Dependencies []string
}

func Remove(pc *config.ProjectConfig, opts RemoveOptions) error {
	if pc.Manifest.Dependencies == nil || len(pc.Manifest.Dependencies) == 0 {
		return fmt.Errorf("cannot remove dependency: there is no dependency in the project")
	}
	for _, dep := range opts.Dependencies {
		if _, ok := pc.Manifest.Dependencies[dep]; !ok {
			fmt.Printf("cannot remove dependency \"%s\": there is no dependency with this name in the project", dep)
		} else {
			delete(pc.Manifest.Dependencies, dep)
		}
	}
	if err := pc.Save(); err != nil {
		return fmt.Errorf("error saving project config: %w", err)
	}
	return nil
}
