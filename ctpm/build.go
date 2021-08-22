package ctpm

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/dependencies"
)

type DependencyFetcher struct{}

func (d DependencyFetcher) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
	libPath := config.LibCachePath(request.Name, request.Version)
	pc, err := config.Load(libPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	ret := make(dependencies.Dependencies)
	for k, v := range pc.Manifest.Dependencies {
		ret[k] = v
	}
	return ret, nil
}

func (d DependencyFetcher) PreAct(_ dependencies.PackageRequest) error  { return nil }
func (d DependencyFetcher) PostAct(_ dependencies.PackageRequest) error { return nil }

func getAllDependencies(pc *config.ProjectConfig) error {
	allDeps := make(manifest.Dependencies)
	for name, version := range pc.Manifest.Dependencies {
		deps, err := dependencies.Install(dependencies.PackageRequest{Name: name, Version: version}, DependencyFetcher{})
		if err != nil {
			return err
		}
		for dname, dversions := range deps {
			for dversion := range dversions {
				allDeps[dname] = dversion
			}
		}
	}
	pc.Manifest.Dependencies = allDeps
	return nil
}

func Build(pc *config.ProjectConfig) error {
	getter := adapter.AdapterGetterImp{}
	adp, err := getter.FromPC(pc.Manifest.Build.Adapter)
	if err != nil {
		return err
	}
	err = getAllDependencies(pc)
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
		if err := Install(dep, semverVersion.String()); err != nil {
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
