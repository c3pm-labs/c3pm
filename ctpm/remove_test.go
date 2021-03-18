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

		BeforeEach(func() {
			var err error
			m := manifest.New()
			m.Name = projectName
			m.Version, err = manifest.VersionFromString("1.0.0")
			Ω(err).ShouldNot(HaveOccurred())

			m.Type = projectType
			m.License = projectLicense
			m.Description = projectDesc
			m.Dependencies = dependencies
			pc = &config.ProjectConfig{Manifest: m, ProjectRoot: projectFolder}
			err = os.MkdirAll(projectFolder, os.ModePerm)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = os.Create(filepath.Join(projectFolder, "c3pm.yml"))
			Ω(err).ShouldNot(HaveOccurred())

		})

		It("remove one dependency", func() {
			toRemove := []string{"calculator"}
			err := ctpm.Remove(pc, ctpm.RemoveOptions{Dependencies: toRemove})
			Ω(err).ShouldNot(HaveOccurred())
			_, ok := pc.Manifest.Dependencies[toRemove[0]]
			Ω(ok).Should(BeFalse())

		})
	})
})
