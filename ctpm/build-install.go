package ctpm

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/registry"
	"os"
	"strings"
)

func getRequiredVersionBuild(dep string, options AddOptions) (name string, version *semver.Version, err error) {
	if err := validateDependency(dep); err != nil {
		return "", nil, err
	}
	dependency := strings.Split(dep, "@")
	if len(dependency) == 1 {
		version, err = registry.GetLastVersion(dep, registry.Options{
			RegistryURL: options.RegistryURL,
		})
		name = dep
		return
	}
	name = dependency[0]
	version, err = semver.NewVersion(dependency[1])
	if err != nil {
		return "", nil, fmt.Errorf("invalid semver string: %w", err)
	}
	return
}

func addDependencyBuild(man *manifest.Manifest, dependency string, version string, opts AddOptions) error {
	options := buildOptions(opts)
	//name, version, err := getRequiredVersionBuild(dependency, options)
	semverVersion, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("error getting dependencies: %w", err)
	}
	var pkg *os.File
	if pkg, err = registry.FetchPackage(dependency, semverVersion, registry.Options{
		RegistryURL: options.RegistryURL,
	}); err != nil {
		return fmt.Errorf("error fetching package: %w", err)
	}
	pkgDir, err := unpackPackage(pkg)
	if err != nil {
		return fmt.Errorf("error unpacking package: %w", err)
	}
	if err = createBuildDirectory(dependency, semverVersion.String()); err != nil {
		return fmt.Errorf("error creating internal c3pm directories: %w", err)
	}
	if err = installPackage(pkgDir); err != nil {
		return fmt.Errorf("error building dependency: %w", err)
	}
	if man.Dependencies == nil {
		man.Dependencies = make(manifest.Dependencies)
	}
	man.Dependencies[dependency] = semverVersion.String()
	return nil
}

func addAllDependencies(pc *config.ProjectConfig) error {
	opts := AddOptions{Force: false, RegistryURL: "", Dependencies: nil}
	for dep, version := range pc.Manifest.Dependencies {
		fmt.Println(dep, version, opts)
		if err := addDependencyBuild(&pc.Manifest, dep, version, opts); err != nil {
			return fmt.Errorf("error adding %s: %w", dep, err)
		}
	}
	return nil
}
