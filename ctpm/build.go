package ctpm

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/config"
)

func Build(pc *config.ProjectConfig) error {
	adp, err := adapter.FromPC(pc)
	if err != nil {
		return err
	}
	return adp.Build(pc)
}

func addAllDependencies(pc *config.ProjectConfig) error {
	for dep, version := range pc.Manifest.Dependencies {
		semverVersion, err := semver.NewVersion(version)
		if err != nil {
			return fmt.Errorf("error getting dependencies: %w", err)
		}
		if err := Install(dep, semverVersion); err != nil {
			return fmt.Errorf("error adding %s: %w", dep, err)
		}
	}
	return nil
}

func AddDependenciesAndBuild(pc *config.ProjectConfig) error {
	err := addAllDependencies(pc)
	if err != nil {
		return fmt.Errorf("error installing dependencies: %w", err)
	}
	err = Build(pc)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	return nil
}
