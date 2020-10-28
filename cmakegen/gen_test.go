package cmakegen

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func getTestFolder(path string) string {
	return filepath.Join("testsArtifacts/", path)
}

func TestGlobbingExprsToCMake(t *testing.T) {
	g := NewGomegaWithT(t)
	cwd, err := os.Getwd()
	g.Expect(err).To(BeNil())
	defer func() {
		err = os.Chdir(cwd)
		g.Expect(err).To(BeNil())
		err = os.RemoveAll(getTestFolder(""))
		g.Expect(err).To(BeNil())
	}()
	testFolder := getTestFolder("glob")
	err = os.MkdirAll(testFolder, os.ModePerm)
	g.Expect(err).To(BeNil())
	err = os.Chdir(testFolder)
	g.Expect(err).To(BeNil())

	err = ioutil.WriteFile("main.c", []byte{}, 0644)
	g.Expect(err).To(BeNil())
	err = ioutil.WriteFile("help.c", []byte{}, 0644)
	g.Expect(err).To(BeNil())

	utilsFolder := "utils"
	err = os.MkdirAll(utilsFolder, os.ModePerm)
	g.Expect(err).To(BeNil())
	err = ioutil.WriteFile(filepath.Join(utilsFolder, "read.c"), []byte{}, 0644)
	g.Expect(err).To(BeNil())
	err = ioutil.WriteFile(filepath.Join(utilsFolder, "write.c"), []byte{}, 0644)
	g.Expect(err).To(BeNil())

	utilsTestsFolder := filepath.Join(utilsFolder, "tests")
	err = os.MkdirAll(utilsTestsFolder, os.ModePerm)
	g.Expect(err).To(BeNil())
	err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_read.c"), []byte{}, 0644)
	g.Expect(err).To(BeNil())
	err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_write.c"), []byte{}, 0644)
	g.Expect(err).To(BeNil())

	files, err := globbingExprToFiles("**/*.c")
	g.Expect(err).To(BeNil())
	g.Expect(files).To(ContainElements([]string{"utils/tests/test_read.c", "utils/tests/test_write.c", "main.c", "help.c", "utils/read.c", "utils/write.c"}))
}
