package manifest

import (
	"fmt"
	"os"

	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCMakeGen(t *testing.T) {
	RegisterFailHandler(Fail)
	path, err := filepath.Abs("testsArtifacts")
	fmt.Println("PATH:", path)
	if err != nil {
		t.Fatal("Failed to get testsArtifacts absolute path")
	}
	RunSpecs(t, "CMakeGenSuite Suite")
	err = os.RemoveAll(path)
	if err != nil {
		t.Fatal("Failed to clean test artifacts\n")
	}
}
