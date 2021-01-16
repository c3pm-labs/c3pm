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
		Ω(err).ShouldNot(HaveOccurred())
		defer func() {
			err = os.Chdir(cwd)
			Ω(err).ShouldNot(HaveOccurred())
			err = os.RemoveAll(getTestFolder(""))
			Ω(err).ShouldNot(HaveOccurred())
		}()
		testFolder := getTestFolder("glob")
		err = os.MkdirAll(testFolder, os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())
		err = os.Chdir(testFolder)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile("main.c", []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile("help.c", []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		utilsFolder := "utils"
		err = os.MkdirAll(utilsFolder, os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile(filepath.Join(utilsFolder, "read.c"), []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile(filepath.Join(utilsFolder, "write.c"), []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		utilsTestsFolder := filepath.Join(utilsFolder, "tests")
		err = os.MkdirAll(utilsTestsFolder, os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_read.c"), []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		err = ioutil.WriteFile(filepath.Join(utilsTestsFolder, "test_write.c"), []byte{}, 0644)
		Ω(err).ShouldNot(HaveOccurred())
		files, err := globbingExprToFiles("**/*.c")
		Ω(err).ShouldNot(HaveOccurred())
		Ω(files).To(ContainElements([]string{"utils/tests/test_read.c", "utils/tests/test_write.c", "main.c", "help.c", "utils/read.c", "utils/write.c"}))
	})
})
