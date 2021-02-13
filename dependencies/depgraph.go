package dependencies

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

func (ps Packages) merge(other Packages) {
	for name, versions := range other {
		for version, _ := range versions {
			ps.add(name, version)
		}
	}
}

type Dependencies map[string]string

type PackageRequest struct {
	Name string
	Version string
}

func Install(r PackageRequest, depHandler DependencyHandler) (Packages, error) {
	installedPackages := Packages{}
	err := depHandler.Install(r)
	if err != nil {
		return nil, err
	}
	installedPackages.add(r.Name, r.Version)
	deps, err := depHandler.FetchDeps(r)
	if err != nil {
		return installedPackages, err
	}
	if deps == nil {
		return installedPackages, nil
	}
	for name, version := range deps {
		dependencies, err := Install(PackageRequest{
			Name:    name,
			Version: version,
		}, depHandler)
		if err != nil {
			return installedPackages, err
		}
		installedPackages.merge(dependencies)
	}
	return installedPackages, nil
}