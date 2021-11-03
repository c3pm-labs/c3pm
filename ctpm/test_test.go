package ctpm_test

import (
	"bytes"
	"github.com/c3pm-labs/c3pm/adapter/defaultadapter"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var testSrc = `#include <iostream>
int main() {
	std::cout << "Test test" << std::endl;
	return 0;
}
`

var _ = Describe("Test", func() {
	Describe("Test executable build", func() {
		projectFolder := getTestFolder("TestTestFolder")
		projectRoot, _ := filepath.Abs(projectFolder)
		m := manifest.New()
		pc := &config.ProjectConfig{Manifest: m, ProjectRoot: projectRoot}
		var wd string

		BeforeEach(func() {
			var err error
			err = os.MkdirAll(projectFolder, os.ModePerm)
			Ω(err).ShouldNot(HaveOccurred())
			wd, err = os.Getwd()
			Ω(err).ShouldNot(HaveOccurred())
			err = os.Chdir(projectFolder)
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			err := os.Chdir(wd)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("Run test", func() {
			pc.Manifest.Name = "buildTest"
			pc.Manifest.Type = manifest.Executable
			pc.Manifest.Build = &manifest.BuildConfig{
				Adapter: &manifest.AdapterConfig{
					Name:    "c3pm",
					Version: defaultadapter.CurrentVersion,
				},
				Config: &defaultadapter.Config{
					Sources:     []string{"main.cpp"},
					TestSources: []string{"test.cpp"},
				},
			}
			pc.Manifest.Version, _ = manifest.VersionFromString("1.0.0")

			err := ioutil.WriteFile("main.cpp", []byte(execSrc), 0644)
			Ω(err).ShouldNot(HaveOccurred())

			err = ioutil.WriteFile("test.cpp", []byte(testSrc), 0644)
			Ω(err).ShouldNot(HaveOccurred())

			err = ctpm.Test(pc)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("Generate cmake scripts", func() {
			fi, err := os.Stat(".c3pm/cmake/CMakeLists.txt")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fi.Size()).To(BeNumerically(">", 0))
		})
		It("Generate makefiles", func() {
			fi, err := os.Stat(".c3pm/build/Makefile")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fi.Size()).To(BeNumerically(">", 0))
		})
		It("ran tests", func() {
			fi, err := os.Stat("buildTest_test")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fi.Size()).To(BeNumerically(">", 0))
			buf := bytes.NewBuffer([]byte{})
			cmd := exec.Command("./buildTest_test")
			cmd.Stdout = buf
			err = cmd.Start()
			Ω(err).ShouldNot(HaveOccurred())

			err = cmd.Wait()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(buf.String()).To(Equal("Test test\n"))
		})
	})
})
