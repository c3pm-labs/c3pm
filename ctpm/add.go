package ctpm

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/env"
	"github.com/c3pm-labs/c3pm/registry"
	"regexp"
	"strings"
)

func Add(pc *config.ProjectConfig, opts AddOptions) error {
	options := buildOptions(opts)
	for _, dep := range opts.Dependencies {
		name, version, err := getRequiredVersion(dep, options)
		if err != nil {
			return fmt.Errorf("error getting dependencies: %w", err)
		}
		err = Install(name, version)
		if err != nil {
			return fmt.Errorf("error adding %s: %w", dep, err)
		}
		if pc.Manifest.Dependencies == nil {
			pc.Manifest.Dependencies = make(manifest.Dependencies)
		}
		pc.Manifest.Dependencies[name] = version.String()
	}
	if err := pc.Save(); err != nil {
		return fmt.Errorf("error saving project config: %w", err)
	}
	return nil
}

type AddOptions struct {
	Force       bool
	RegistryURL string

	Dependencies []string
}

func buildOptions(opts AddOptions) AddOptions {
	if opts.RegistryURL == "" {
		opts.RegistryURL = env.REGISTRY_ENDPOINT
	}
	return opts
}

const depRegexString = `^[\-a-zA-Z0-9_]*(@.*)?$`

var depRegex *regexp.Regexp

func getDepRegexp() (regex *regexp.Regexp, err error) {
	if depRegex == nil {
		depRegex, err = regexp.Compile(depRegexString)
	}
	regex = depRegex
	return
}

func validateDependency(dep string) error {
	regex, err := getDepRegexp()
	if err != nil {
		return fmt.Errorf("error compiling validation regexp: %w", err)
	}
	if regex.MatchString(dep) {
		return nil
	}
	return errors.New("%s is not a valid dependency string")
}

func getRequiredVersion(dep string, options AddOptions) (name string, version *semver.Version, err error) {
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
