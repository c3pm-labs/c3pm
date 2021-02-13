package dependencies_test

import (
	"errors"
	"fmt"
	"github.com/c3pm-labs/c3pm/dependencies"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PackageMeta struct {
	Version string
	Dependencies dependencies.Dependencies
}

var packages = map[string][]PackageMeta{
	"boost/assert": {{
		Version:      "0.0.0",
		Dependencies: dependencies.Dependencies{"boost/config": "0.0.0"},
	}, {
		Version:      "0.0.1",
		Dependencies: dependencies.Dependencies{"boot/config": "0.0.0"},
	}},
	"boost/config": {{"0.0.0", nil}},
}

type TestFetcher int

func (TestFetcher) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
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

func (TestFetcher) Install(request dependencies.PackageRequest) error {
	return nil
}

const fetcher TestFetcher = 0

var _ = Describe("Dependencies", func() {
	It("fetches no dependencies if the package has none", func() {
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/config",
			Version: "0.0.0",
		}, fetcher)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).ShouldNot(BeNil())
	})

	It("fetches the correct dependencies with a single dependency", func() {
		pkgs, err := dependencies.Install(dependencies.PackageRequest{
			Name:    "boost/assert",
			Version: "0.0.0",
		}, fetcher)
		fmt.Println(pkgs)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(pkgs).ShouldNot(BeNil())
	})
})
