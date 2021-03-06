package defaultadapter

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Gen Test", func() {
	path, _ := filepath.Abs("../../test_helpers/projects/genProject")
	var (
		simpleProject = &config.ProjectConfig{
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
				Build: &manifest.BuildConfig{
					Adapter: &manifest.AdapterConfig{
						Name:    "c3pm",
						Version: CurrentVersion,
					},
					Config: &Config{
						Sources:     []string{"**/*.cpp"},
						Headers:     []string{"**/*.hpp"},
						IncludeDirs: []string{"include"},
					},
				},
				Dependencies: manifest.Dependencies{},
			},
			ProjectRoot: path,
		}
		projectWithDependencies = &config.ProjectConfig{
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
				Build: &manifest.BuildConfig{
					Adapter: &manifest.AdapterConfig{
						Name:    "c3pm",
						Version: CurrentVersion,
					},
					Config: &Config{
						Sources:     []string{"**/*.cpp"},
						Headers:     []string{"**/*.hpp"},
						IncludeDirs: []string{"include"},
					},
				},
				Dependencies: manifest.Dependencies{
					"hello": "1.0.3",
					"m":     "2.0.0",
				},
			},
			ProjectRoot: path,
		}
	)

	BeforeEach(func() {
	})
	AfterEach(func() {
		err := os.RemoveAll(cmakeDirFromPc(simpleProject))
		Ω(err).ShouldNot(HaveOccurred())
	})
	Context("generates a cmake file without dependencies", func() {
		err := generateCMakeScripts(cmakeDirFromPc(simpleProject), simpleProject, nil)
		fmt.Println(err)
		Ω(err).ShouldNot(HaveOccurred())
		data, err := ioutil.ReadFile(filepath.Join(cmakeDirFromPc(simpleProject), "CMakeLists.txt"))
		Ω(err).ShouldNot(HaveOccurred())
		content := string(data)

		It("contains the correct source files", func() {
			mainPath, err := filepath.Abs(filepath.Join(path, "main.cpp"))
			Ω(err).ShouldNot(HaveOccurred())
			libPath, err := filepath.Abs(filepath.Join(path, "lib", "hello.cpp"))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(content).Should(ContainSubstring(mainPath))
			Ω(content).Should(ContainSubstring(libPath))
		})
		It("doesn't contain dependencies", func() {
			Ω(content).ShouldNot(ContainSubstring("-l"))
			Ω(content).ShouldNot(ContainSubstring("-L"))
		})
	})
	Context("generates a cmake file with dependencies", func() {
		_ = projectWithDependencies
		//TODO: dependencies tests
		//err := cmakegen.generateCMakeScripts(projectWithDependencies)
		//Ω(err).ShouldNot(HaveOccurred())
		//data, err := ioutil.ReadFile(filepath.Join(projectWithDependencies.CMakeDir(), "CMakeLists.txt"))
		//Ω(err).ShouldNot(HaveOccurred())
		//content := string(data)
		//fmt.Println(content)
		//It("contains the correct source files", func() {
		//	mainPath, err := filepath.Abs(filepath.Join(path, "main.cpp"))
		//	Ω(err).ShouldNot(HaveOccurred())
		//	libPath, err := filepath.Abs(filepath.Join(path, "lib", "hello.cpp"))
		//	Ω(err).ShouldNot(HaveOccurred())
		//	Ω(content).Should(ContainSubstring(mainPath))
		//	Ω(content).Should(ContainSubstring(libPath))
		//})
		//It("contains links to the dependencies", func() {
		//	Ω(content).Should(ContainSubstring("-lhello"))
		//	Ω(content).Should(ContainSubstring("-lm"))
		//	Ω(content).Should(ContainSubstring("-L"))
		//})
	})
})
