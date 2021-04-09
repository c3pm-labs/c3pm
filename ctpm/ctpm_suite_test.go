package ctpm_test

import (
	"os"

	//"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCtpm(t *testing.T) {
	RegisterFailHandler(Fail)
	path, err := filepath.Abs("testsArtifacts")
	if err != nil {
		t.Fatal("Failed to get testsArtifacts absolute path")
	}
	RunSpecs(t, "Ctpm Suite")
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal("Failed to clean test artifacts\n")
	}
}

func getTestFolder(path string) string {
	return filepath.Join("testsArtifacts/", path)
}
