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

var _ = Describe("Init", func() {
	Describe("Project creation", func() {
		var projectFolder = getTestFolder("InitTestFolder")
		projectFolder, err := filepath.Abs(projectFolder)
		Ω(err).ShouldNot(HaveOccurred())

		var projectName = "InitProject"
		var projectType = manifest.Library
		var projectDesc = "description"
		var projectLicense = "MIT"

		When("initializing manifest", func() {
			It("Did not fail", func() {
				var err error
				m := manifest.New()
				m.Name = projectName
				m.Version, err = manifest.VersionFromString("1.0.0")
				Ω(err).ShouldNot(HaveOccurred())

				m.Type = projectType
				m.License = projectLicense
				m.Description = projectDesc
				pc := &config.ProjectConfig{Manifest: m, ProjectRoot: projectFolder}
				err = ctpm.Init(pc, ctpm.InitDefaultOptions)
				Ω(err).ShouldNot(HaveOccurred())
			})
		})

		It("Create c3pm.yml file", func() {
			pc, err := manifest.Load(filepath.Join(projectFolder, "c3pm.yml"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(pc.Name).To(Equal(projectName))
			Ω(pc.License).To(Equal(projectLicense))
			Ω(pc.Type).To(Equal(projectType))
			Ω(pc.Description).To(Equal(projectDesc))
		})
		// 		PIt("Project template of either library or an executable")
		It("Depending on the user's choice, a LICENSE file", func() {
			lic, err := os.Stat(filepath.Join(projectFolder, "LICENSE"))
			Ω(err).ShouldNot(HaveOccurred())

			Ω(lic).ShouldNot(BeNil())
		})
	})
})
