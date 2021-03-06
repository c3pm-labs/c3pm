package ctpm_test

import (
	"bytes"
	"github.com/c3pm-labs/c3pm/adapter/defaultadapter"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var execSrc = `#include <iostream>
int main() {
	std::cout << "Build test" << std::endl;
}
`

var _ = Describe("Build", func() {
	Describe("Executable build", func() {
		projectFolder := getTestFolder("BuildTestFolder")
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
		It("Run build", func() {
			pc.Manifest.Name = "buildTest"
			pc.Manifest.Type = manifest.Executable
			pc.Manifest.Build = &manifest.BuildConfig{
				Adapter: &manifest.AdapterConfig{
					Name:    "c3pm",
					Version: defaultadapter.CurrentVersion,
				},
				Config: &defaultadapter.Config{
					Sources: []string{"main.cpp"},
				},
			}
			pc.Manifest.Version, _ = manifest.VersionFromString("1.0.0")

			err := ioutil.WriteFile("main.cpp", []byte(execSrc), 0644)
			Ω(err).ShouldNot(HaveOccurred())

			err = ctpm.Build(pc)
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
		It("build m", func() {
			fi, err := os.Stat("buildTest")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fi.Size()).To(BeNumerically(">", 0))
			buf := bytes.NewBuffer([]byte{})
			cmd := exec.Command("./buildTest")
			cmd.Stdout = buf
			err = cmd.Start()
			Ω(err).ShouldNot(HaveOccurred())

			err = cmd.Wait()
			Ω(err).ShouldNot(HaveOccurred())

			Ω(buf.String()).To(Equal("Build test\n"))
		})
	})
})
