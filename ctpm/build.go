package ctpm

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/adapter_interface"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/dependencies"
)

type DependencyBuilder struct {
	Done manifest.Dependencies
}

func (d DependencyBuilder) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
	if _, ok := d.Done[request.Name]; ok {
		return dependencies.Dependencies{}, nil
	}
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

func (d DependencyBuilder) PreAct(_ dependencies.PackageRequest) error { return nil }
func (d DependencyBuilder) PostAct(request dependencies.PackageRequest) error {
	fmt.Printf("Building %s:%s\n", request.Name, request.Version)

	libPath := config.LibCachePath(request.Name, request.Version)
	pc, err := config.Load(libPath)
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	getter := adapter.AdapterGetterImp{}
	var adp adapter_interface.Adapter
	adp, err = getter.FromPC(pc.Manifest.Build.Adapter)
	if err != nil {
		return err
	}
	err = adp.Build(pc)
	if err != nil {
		return fmt.Errorf("error building: %w", err)
	}
	d.Done[fmt.Sprintf(request.Name)] = request.Version
	return nil
}

func getAllDependencies(pc *config.ProjectConfig) error {
	allDeps := make(manifest.Dependencies)
	allDeps[pc.Manifest.Name] = pc.Manifest.Version.String()
	for name, version := range pc.Manifest.Dependencies {
		_, err := dependencies.Install(dependencies.PackageRequest{Name: name, Version: version}, DependencyBuilder{Done: allDeps})
		if err != nil {
			return err
		}
	}
	delete(allDeps, pc.Manifest.Name)
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
