package cmakegen

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gen Test", func() {
	It("generate cmake", func() {
		cwd, err := os.Getwd()
		Ω(err).To(BeNil())
		defer func() {
			err = os.Chdir(cwd)
			Ω(err).To(BeNil())
			err = os.RemoveAll(getTestFolder(""))
			Ω(err).To(BeNil())
		}()
		testFolder := getTestFolder("glob")
		err = os.MkdirAll(testFolder, os.ModePerm)
		Ω(err).To(BeNil())
		err = os.Chdir(testFolder)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile("main.c", []byte{}, 0644)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile("help.c", []byte{}, 0644)
		Ω(err).To(BeNil())
		utilsFolder := "utils"
		err = os.MkdirAll(utilsFolder, os.ModePerm)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile(filepath.Join(utilsFolder, "read.c"), []byte{}, 0644)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile(filepath.Join(utilsFolder, "write.c"), []byte{}, 0644)
		Ω(err).To(BeNil())
		utilsTestsFolder := filepath.Join(utilsFolder, "tests")
		err = os.MkdirAll(utilsTestsFolder, os.ModePerm)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_read.c"), []byte{}, 0644)
		Ω(err).To(BeNil())
		err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_write.c"), []byte{}, 0644)
		Ω(err).To(BeNil())
		files, err := globbingExprToFiles("**/*.c")
		Ω(err).To(BeNil())
		Ω(files).To(ContainElements([]string{"utils/tests/test_read.c", "utils/tests/test_write.c", "main.c", "help.c", "utils/read.c", "utils/write.c"}))
	})
})
