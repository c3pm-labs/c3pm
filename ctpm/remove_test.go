package ctpm_test

import (
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Remove", func() {
	Describe("Remove dependencies", func() {
		var projectFolder = getTestFolder("RemoveTestFolder")
		var projectName = "Test-Remove"
		var projectType = manifest.Executable
		var projectDesc = "description"
		var projectLicense = "MIT"
		var dependencies = map[string]string{
			"calculator":  "1.0.0",
			"sort":        "1.0.0",
			"simple-math": "1.0.0",
			"super-math":  "1.0.0",
		}
		var pc *config.ProjectConfig
		err := os.MkdirAll(projectFolder, os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())

		BeforeEach(func() {
			m := manifest.New()
			m.Name = projectName
			m.Version, err = manifest.VersionFromString("1.0.0")
			Ω(err).ShouldNot(HaveOccurred())
			m.Type = projectType
			m.License = projectLicense
			m.Description = projectDesc
			m.Dependencies = dependencies

			pc = &config.ProjectConfig{Manifest: m, ProjectRoot: projectFolder}

			_, err = os.Create(filepath.Join(projectFolder, "c3pm.yml"))
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.Remove(filepath.Join(projectFolder, "c3pm.yml"))
		})

		It("remove one dependency", func() {
			toRemove := []string{"calculator"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).ShouldNot(HaveOccurred())
			_, ok := pc.Manifest.Dependencies[toRemove[0]]
			Ω(ok).Should(BeFalse())

		})
		It("remove several dependencies at the same time", func() {
			toRemove := []string{"calculator", "sort", "simple-math"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).ShouldNot(HaveOccurred())
			for _, dep := range toRemove {
				_, ok := pc.Manifest.Dependencies[dep]
				Ω(ok).Should(BeFalse())
			}
		})
		It("remove unexisting dependency", func() {
			toRemove := []string{"toto"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).ShouldNot(HaveOccurred())
			for _, dep := range toRemove {
				_, ok := pc.Manifest.Dependencies[dep]
				Ω(ok).Should(BeFalse())
			}
		})
		It("remove dependency from an empty dependency array", func() {
			pc.Manifest.Dependencies = map[string]string{}
			toRemove := []string{"toto"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).Should(HaveOccurred())
		})
		It("remove dependency from an nil dependency array", func() {
			pc.Manifest.Dependencies = nil
			toRemove := []string{"toto"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).Should(HaveOccurred())
		})
	})
})
