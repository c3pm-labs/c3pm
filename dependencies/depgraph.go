package dependencies

import (
	"fmt"
	"strings"
)

// Packages is a map from names to the list of versions.
// The versions are represented as a map of string over struct{} to use the speed and uniqueness of map keys.
type Packages map[string]map[string]struct{}

func (ps Packages) Find(name, version string) bool {
	vs, ok := ps[name]
	if !ok {
		return false
	}
	_, ok = vs[version]
	return ok
}

func (ps Packages) add(name, version string) {
	vs, ok := ps[name]
	if !ok {
		ps[name] = map[string]struct{}{version: {}}
	} else {
		vs[version] = struct{}{}
	}
}

func (ps Packages) String() string {
	sb := strings.Builder{}
	i := 0
	for name, versions := range ps {
		sb.WriteString(name + "@")
		sb.WriteString("[")
		j := 0
		for version := range versions {
			sb.WriteString(version)
			if j != len(versions)-1 && len(versions) >= 1 {
				sb.WriteString(",")
			}
			j += 1
		}
		sb.WriteString("]")
		if i != len(ps)-1 && len(ps) >= 1 {
			sb.WriteString(", ")
		}
		i += 1
	}
	return sb.String()
}

type Dependencies map[string]string

type PackageRequest struct {
	Name    string
	Version string
}

func install(r PackageRequest, depHandler DependencyHandler, installedPackages *Packages) error {
	// If package has already been installed, pass
	if installedPackages.Find(r.Name, r.Version) {
		return nil
	}
	err := depHandler.PreAct(r)
	if err != nil {
		return fmt.Errorf("error during prefetch action for %s/%s: %w", r.Name, r.Version, err)
	}
	installedPackages.add(r.Name, r.Version)
	deps, err := depHandler.FetchDeps(r)
	if err != nil {
		return fmt.Errorf("error fetching dependencies of %s/%s: %w", r.Name, r.Version, err)
	}
	if deps == nil {
		return nil
	}
	for name, version := range deps {
		err := install(PackageRequest{
			Name:    name,
			Version: version,
		}, depHandler, installedPackages)
		if err != nil {
			return err
		}
	}
	err = depHandler.PostAct(r)
	if err != nil {
		return fmt.Errorf("error during postfetch action for %s/%s: %w", r.Name, r.Version, err)
	}
	return nil
}

func Install(r PackageRequest, depHandler DependencyHandler) (Packages, error) {
	installedPackages := Packages{}
	err := install(r, depHandler, &installedPackages)
	return installedPackages, err
}
