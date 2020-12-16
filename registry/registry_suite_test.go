package registry

import (
	"fmt"
	"os"

	//"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmakeGen(t *testing.T) {
	RegisterFailHandler(Fail)
	path, err := filepath.Abs("testsArtifacts")
	fmt.Println("PATH:", path)
	if err != nil {
		t.Fatal("Failed to get testsArtifacts absolute path")
	}
	RunSpecs(t, "CmakeGenSuite Suite")
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal("Failed to clean test artifacts\n")
	}
}

func getTestFolder(path string) string {
	return filepath.Join("testsArtifacts/", path)
}
