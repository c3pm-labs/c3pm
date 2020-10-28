package ctpm_test

import (
	"github.com/gabrielcolson/c3pm/cli/config"
	"github.com/gabrielcolson/c3pm/cli/config/manifest"
	"github.com/gabrielcolson/c3pm/cli/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Init", func() {
	Describe("Project creation", func() {
		var projectFolder = getTestFolder("InitTestFolder")
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
				Ω(err).To(BeNil())
				m.Type = projectType
				m.License = projectLicense
				m.Description = projectDesc
				pc := &config.ProjectConfig{Manifest: m, ProjectRoot: projectFolder}
				err = ctpm.Init(pc, ctpm.InitDefaultOptions)
				Ω(err).To(BeNil())
			})
		})

		It("Create c3pm.yml file", func() {
			pc, err := manifest.Load(filepath.Join(projectFolder, "c3pm.yml"))
			Ω(err).To(BeNil())
			Ω(pc.Name).To(Equal(projectName))
			Ω(pc.License).To(Equal(projectLicense))
			Ω(pc.Type).To(Equal(projectType))
			Ω(pc.Description).To(Equal(projectDesc))
		})
		// 		PIt("Project template of either library or an executable")
		It("A .c3pm directory", func() {
			fileInfo, err := os.Stat(filepath.Join(projectFolder, ".c3pm"))
			Ω(err).To(BeNil())
			Ω(fileInfo).ShouldNot(BeNil())
		})
		It("A CMakeLists.txt", func() {
			cmake, err := os.Stat(filepath.Join(projectFolder, ".c3pm/cmake/CMakeLists.txt"))
			Ω(err).To(BeNil())
			Ω(cmake).ShouldNot(BeNil())
		})
		It("Depending on the user's choice, a LICENSE file", func() {
			lic, err := os.Stat(filepath.Join(projectFolder, "LICENSE"))
			Ω(err).To(BeNil())
			Ω(lic).ShouldNot(BeNil())
		})
	})
})
