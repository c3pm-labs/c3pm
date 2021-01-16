package ctpm_test

import (
	"bytes"
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
		var oldPath string
		projectFolder := getTestFolder("BuildTestFolder")
		projectRoot, _ := filepath.Abs(projectFolder)
		m := manifest.New()
		pc := &config.ProjectConfig{Manifest: m, ProjectRoot: projectRoot}

		It("Run build", func() {
			pc.Manifest.Name = "buildTest"
			pc.Manifest.Type = manifest.Executable
			pc.Manifest.Files.Sources = []string{"main.cpp"}
			pc.Manifest.Version, _ = manifest.VersionFromString("1.0.0")

			err := os.MkdirAll(projectFolder, os.ModePerm)
			Ω(err).To(BeNil())
			oldPath, err = filepath.Rel(projectFolder, ".")
			Ω(err).To(BeNil())
			err = os.Chdir(projectFolder)
			Ω(err).To(BeNil())
			err = ioutil.WriteFile("main.cpp", []byte(execSrc), 0644)
			Ω(err).To(BeNil())
			err = ctpm.Build(pc)
		})

		It("Generate cmake scripts", func() {
			fi, err := os.Stat(".c3pm/cmake/CMakeLists.txt")
			Ω(err).To(BeNil())
			Ω(fi.Size()).To(BeNumerically(">", 0))
		})
		It("Generate makefiles", func() {
			fi, err := os.Stat(".c3pm/build/Makefile")
			Ω(err).To(BeNil())
			Ω(fi.Size()).To(BeNumerically(">", 0))
		})
		It("build m", func() {
			fi, err := os.Stat("buildTest")
			Ω(err).To(BeNil())
			Ω(fi.Size()).To(BeNumerically(">", 0))
			buf := bytes.NewBuffer([]byte{})
			cmd := exec.Command("./buildTest")
			cmd.Stdout = buf
			err = cmd.Start()
			Ω(err).To(BeNil())
			err = cmd.Wait()
			Ω(err).To(BeNil())
			Ω(buf.String()).To(Equal("Build test\n"))
		})

		_ = os.Chdir(oldPath)
	})
})
