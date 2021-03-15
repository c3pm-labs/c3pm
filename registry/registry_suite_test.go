package registry

import (
	"os"

	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCMakeGen(t *testing.T) {
	RegisterFailHandler(Fail)
	path, err := filepath.Abs("testsArtifacts")
	if err != nil {
		t.Fatal("Failed to get testsArtifacts absolute path")
	}
	RunSpecs(t, "CMakeGenSuite Suite")
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal("Failed to clean test artifacts\n")
	}
}
