package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/adapter/defaultadapter"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/mohae/deepcopy"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	OriginalDir = "../test_helpers/yamls"
)

var (
	OriginalDirAbs, _ = filepath.Abs(OriginalDir)
	TestConfig        = &config.ProjectConfig{
		Manifest: manifest.Manifest{
			C3PMVersion: manifest.C3PMVersion1,
			Type:        manifest.Executable,
			Name:        "hello-bin",
			Description: "Demo binary",
			Version: manifest.Version{
				Version: semver.MustParse("1.1.5"),
			},
			Standard: "20",
			License:  "ISC",
			Publish: &manifest.PublishConfig{
				Include: []string{"test-include"},
				Exclude: []string{"test-exclude"},
			},
			Build: &manifest.BuildConfig{
				Adapter: &manifest.AdapterConfig{
					Name:    "c3pm",
					Version: defaultadapter.CurrentVersion,
				},
				Config: nil,
			},
			Documentation: "",
			Website:       "",
			Repository:    "",
			Contributors:  "",
			Dependencies: manifest.Dependencies{
				"hello": "1.0.5",
			},
			TestDependencies: manifest.Dependencies{},
		},
		ProjectRoot: OriginalDirAbs,
	}
)

var _ = Describe("Config loading and writing", func() {

	Context("file reads", func() {
		It("loads the file properly", func() {
			p, err := config.Load(OriginalDir)
			Ω(err).ShouldNot(HaveOccurred())
			p.Manifest.Build.Config = nil // we ignore the Build.Config field
			Ω(p).Should(Equal(TestConfig))
		})
	})

	Context("file writes", func() {
		var (
			TargetDir string
		)

		BeforeEach(func() {
			dir, err := ioutil.TempDir("", "c3pm_test_*")
			Ω(err).ShouldNot(HaveOccurred())
			TargetDir = dir
		})

		AfterEach(func() {
			_ = os.RemoveAll(TargetDir)
		})

		It("Updates the file correctly", func() {
			p := deepcopy.Copy(TestConfig).(*config.ProjectConfig)
			p.ProjectRoot = TargetDir
			p.Manifest.Version = TestConfig.Manifest.Version                             // Private values are not copied by deep copy, so let's just add the value ourselves
			p.Manifest.Build.Adapter.Version = TestConfig.Manifest.Build.Adapter.Version // Private values are not copied by deep copy, so let's just add the value ourselves
			p.Manifest.Description = "Different Description"
			p.Manifest.Build.Config = nil // we ignore the Build.Config field
			err := p.Save()
			Ω(err).ShouldNot(HaveOccurred())
			p2, err := config.Load(TargetDir)
			Ω(err).ShouldNot(HaveOccurred())
			p2.Manifest.Build.Config = nil // we ignore the Build.Config field
			Ω(p2).ShouldNot(Equal(TestConfig), "Test against the original config")
			Ω(p2).Should(Equal(p))
		})
	})
})

var _ = Describe("Config utils", func() {
	Context("Global directory", func() {
		var (
			OriginalHomeDir = os.Getenv("HOME")
		)
		AfterEach(func() {
			err := os.Setenv("HOME", OriginalHomeDir)
			Ω(err).ShouldNot(HaveOccurred())
			err = os.Unsetenv("C3PM_USER_DIR")
			Ω(err).ShouldNot(HaveOccurred())
		})

		Describe("Root global directory", func() {
			It("gets from HOME", func() {
				err := os.Setenv("HOME", ".")
				Ω(err).ShouldNot(HaveOccurred())
				path := config.GlobalC3PMDirPath()
				Ω(path).Should(Equal(".c3pm"))
			})
			It("ets from C3PM_USER_DIR", func() {
				err := os.Setenv("C3PM_USER_DIR", "../.c3pm")
				Ω(err).ShouldNot(HaveOccurred())
				path := config.GlobalC3PMDirPath()
				Ω(path).Should(Equal("../.c3pm"))
			})
			It("has priority from C3PM_USER_DIR over HOME", func() {
				err := os.Setenv("C3PM_USER_DIR", "../.c3pm")
				Ω(err).ShouldNot(HaveOccurred())
				err = os.Setenv("HOME", "..")
				Ω(err).ShouldNot(HaveOccurred())
				path := config.GlobalC3PMDirPath()
				Ω(path).Should(Equal("../.c3pm"))
			})
		})
		Describe("Cache directory", func() {
			It("gets the cache directory", func() {
				err := os.Setenv("C3PM_USER_DIR", ".c3pm")
				Ω(err).ShouldNot(HaveOccurred())
				path := config.LibCachePath("a", "b")
				Ω(path).Should(Equal(".c3pm/cache/a/b"))
			})
		})
	})
})
