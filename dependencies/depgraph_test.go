package dependencies_test

import (
	"errors"
	"github.com/c3pm-labs/c3pm/dependencies"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PackageMeta struct {
	Version      string
	Dependencies dependencies.Dependencies
}

var packages = map[string][]PackageMeta{
	"boost/assert": {{
		Version:      "0.0.0",
		Dependencies: dependencies.Dependencies{"boost/config": "0.0.0"},
	}, {
		Version:      "0.0.1",
		Dependencies: dependencies.Dependencies{"boost/config": "0.0.0"},
	}},
	"boost/chrono":          {{"0.0.0", dependencies.Dependencies{"boost/assert": "0.0.0", "boost/config": "0.0.0", "boost/core": "0.0.0", "boost/mpl": "0.0.0", "boost/static_assert": "0.0.0", "boost/integer": "0.0.0", "boost/typeof": "0.0.0", "boost/type_traits": "0.0.0", "boost/winapi": "0.0.0", "boost/move": "0.0.0", "boost/system": "0.0.0", "boost/throw_exception": "0.0.0", "boost/ratio": "0.0.0"}}},
	"boost/config":          {{"0.0.0", nil}},
	"boost/core":            {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0", "boost/assert": "0.0.0"}}},
	"boost/integer":         {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0", "boost/static_assert": "0.0.0"}}},
	"boost/move":            {{"0.0.0", dependencies.Dependencies{"boost/assert": "0.0.0", "boost/config": "0.0.0", "boost/core": "0.0.0", "boost/static_assert": "0.0.0"}}},
	"boost/mpl":             {{"0.0.0", dependencies.Dependencies{"boost/core": "0.0.0", "boost/static_assert": "0.0.0", "boost/predef": "0.0.0", "boost/preprocessor": "0.0.0", "boost/throw_exception": "0.0.0", "boost/type_traits": "0.0.0", "boost/utility": "0.0.0"}}},
	"boost/predef":          {{"0.0.0", nil}},
	"boost/preprocessor":    {{"0.0.0", nil}},
	"boost/ratio":           {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0", "boost/core": "0.0.0", "boost/mpl": "0.0.0", "boost/integer": "0.0.0", "boost/type_traits": "0.0.0"}}},
	"boost/static_assert":   {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0"}}},
	"boost/system":          {{"0.0.0", dependencies.Dependencies{"boost/core": "0.0.0", "boost/assert": "0.0.0", "boost/predef": "0.0.0"}}},
	"boost/throw_exception": {{"0.0.0", dependencies.Dependencies{"boost/assert": "0.0.0", "boost/config": "0.0.0"}}},
	"boost/type_traits":     {{"0.0.0", dependencies.Dependencies{"boost/core": "0.0.0", "boost/config": "0.0.0", "boost/static_assert": "0.0.0"}}},
	"boost/typeof":          {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0", "boost/core": "0.0.0", "boost/mpl": "0.0.0", "boost/preprocessor": "0.0.0", "boost/type_traits": "0.0.0"}}},
	"boost/utility":         {{"0.0.0", dependencies.Dependencies{"boost/mpl": "0.0.0", "boost/throw_exception": "0.0.0"}}},
	"boost/winapi":          {{"0.0.0", dependencies.Dependencies{"boost/config": "0.0.0", "boost/predef": "0.0.0"}}},
}

type TestDependencyHandler int

func (TestDependencyHandler) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
	versions, ok := packages[request.Name]
	if !ok {
		return nil, errors.New("package not found")
	}
	for _, p := range versions {
		if p.Version == request.Version {
			return p.Dependencies, nil
		}
	}
	return nil, errors.New("version not found")
}

func (TestDependencyHandler) Act(request dependencies.PackageRequest) error {
	return nil
}

const fetcher TestDependencyHandler = 0

var _ = Describe("Dependencies", func() {
	It("fetches no dependencies if the package has none", func() {
		packages := dependencies.Packages{
			"boost/config": {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/config",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
	})

	It("fetches the correct dependencies with a single dependency", func() {
		packages := dependencies.Packages{
			"boost/config": {"0.0.0": struct{}{}},
			"boost/assert": {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/assert",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
	})

	It("fetches the correct dependencies with duplicate dependencies", func() {
		packages := dependencies.Packages{
			"boost/assert": {"0.0.0": struct{}{}},
			"boost/config": {"0.0.0": struct{}{}},
			"boost/core":   {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/core",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
	})

	It("fetches the correct dependencies with two levels of dependency", func() {
		packages := dependencies.Packages{
			"boost/assert": {"0.0.0": struct{}{}},
			"boost/config": {"0.0.0": struct{}{}},
			"boost/core":   {"0.0.0": struct{}{}},
			"boost/predef": {"0.0.0": struct{}{}},
			"boost/system": {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/system",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
	})

	It("resolves circular dependencies", func(done Done) {
		packages := dependencies.Packages{
			"boost/assert":          {"0.0.0": struct{}{}},
			"boost/config":          {"0.0.0": struct{}{}},
			"boost/core":            {"0.0.0": struct{}{}},
			"boost/mpl":             {"0.0.0": struct{}{}},
			"boost/predef":          {"0.0.0": struct{}{}},
			"boost/preprocessor":    {"0.0.0": struct{}{}},
			"boost/type_traits":     {"0.0.0": struct{}{}},
			"boost/throw_exception": {"0.0.0": struct{}{}},
			"boost/static_assert":   {"0.0.0": struct{}{}},
			"boost/utility":         {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/mpl",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
		close(done)
	})

	It("resolves an important (17) number of packages", func(done Done) {
		packages := dependencies.Packages{
			"boost/assert":          {"0.0.0": struct{}{}},
			"boost/config":          {"0.0.0": struct{}{}},
			"boost/core":            {"0.0.0": struct{}{}},
			"boost/chrono":          {"0.0.0": struct{}{}},
			"boost/integer":         {"0.0.0": struct{}{}},
			"boost/move":            {"0.0.0": struct{}{}},
			"boost/mpl":             {"0.0.0": struct{}{}},
			"boost/predef":          {"0.0.0": struct{}{}},
			"boost/preprocessor":    {"0.0.0": struct{}{}},
			"boost/ratio":           {"0.0.0": struct{}{}},
			"boost/type_traits":     {"0.0.0": struct{}{}},
			"boost/typeof":          {"0.0.0": struct{}{}},
			"boost/throw_exception": {"0.0.0": struct{}{}},
			"boost/static_assert":   {"0.0.0": struct{}{}},
			"boost/system":          {"0.0.0": struct{}{}},
			"boost/utility":         {"0.0.0": struct{}{}},
			"boost/winapi":          {"0.0.0": struct{}{}},
		}
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/chrono",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).Should(Equal(packages))
		close(done)

	})
})
