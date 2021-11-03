package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
)

func List(pc *config.ProjectConfig) error {
	err := getAllDependencies(pc)
	if err != nil {
		return err
	}

	for key, val := range pc.Manifest.Dependencies {
		fmt.Printf("%s@%s\n", key, val)
	}

	return nil
}
