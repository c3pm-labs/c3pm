package ctpm

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/cmake"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/env"
	"github.com/c3pm-labs/c3pm/registry"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Add(pc *config.ProjectConfig, opts AddOptions) error {
	for _, dep := range opts.Dependencies {
		if err := addDependency(&pc.Manifest, dep, opts); err != nil {
			return fmt.Errorf("error adding %s: %w", dep, err)
		}
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

func addDependency(man *manifest.Manifest, dependency string, opts AddOptions) error {
	options := buildOptions(opts)
	name, version, err := getRequiredVersion(dependency, options)
	if err != nil {
		return fmt.Errorf("error getting dependencies: %w", err)
	}
	var pkg *os.File
	if pkg, err = registry.FetchPackage(name, version, registry.Options{
		RegistryURL: options.RegistryURL,
	}); err != nil {
		return fmt.Errorf("error fetching package: %w", err)
	}
	pkgDir, err := unpackPackage(pkg)
	if err != nil {
		return fmt.Errorf("error unpacking package: %w", err)
	}
	if err = installPackage(pkgDir); err != nil {
		return fmt.Errorf("error building dependency: %w", err)
	}
	if man.Dependencies == nil {
		man.Dependencies = make(manifest.Dependencies)
	}
	man.Dependencies[name] = version.String()
	return nil
}

func unpackPackage(pkg *os.File) (string, error) {
	tr := tar.NewReader(pkg)
	dst, err := ioutil.TempDir("", "c3pm")
	if err != nil {
		return "", fmt.Errorf("error creating package directory: %w", err)
	}
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			_ = os.Remove(pkg.Name())
			return dst, nil
		case err != nil:
			return "", err
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.ModePerm); err != nil {
					return "", err
				}
			}
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Dir(target), os.ModePerm)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return "", err
			}
			_ = f.Close()
		case tar.TypeSymlink:
			err = os.MkdirAll(filepath.Dir(target), os.ModePerm)
			if err != nil {
				return "", err
			}
			if err = os.Symlink(header.Linkname, target); err != nil {
				return "", err
			}
		}
	}
}

func installPackage(pkgDir string) error {
	pc, err := config.Load(pkgDir)
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	err = Build(pc)
	if err != nil {
		return err
	}
	err = cmake.Install(pc.BuildDir())
	if err != nil {
		return err
	}
	return pc.Manifest.Save(filepath.Join(config.LibCachePath(pc.Manifest.Name, pc.Manifest.Version.String()), "c3pm.yml"))
}

const depRegexString = `^[\-a-z0-9_]*(@.*)?$`

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
