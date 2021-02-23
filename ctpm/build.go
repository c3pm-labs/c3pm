package ctpm

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config"
)

func Build(pc *config.ProjectConfig) error {
	adp, err := adapterFromPc(pc)
	if err != nil {
		return err
	}
	return adp.Build(pc)
}

func addAllDependencies(pc *config.ProjectConfig) error {
	opts := AddOptions{Force: false, RegistryURL: "", Dependencies: nil}
	options := buildOptions(opts)

	for dep, version := range pc.Manifest.Dependencies {
		semverVersion, err := semver.NewVersion(version)
		if err != nil {
			return fmt.Errorf("error getting dependencies: %w", err)
		}
		if err := addDependency(&pc.Manifest, dep, semverVersion, options); err != nil {
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
