package ctpm_test

import (
	"archive/tar"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
)

var _ = Describe("Publish", func() {
	var wd string
	var filesFound []string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		Ω(req.URL.String()).To(Equal("/packages/publish"))
		Ω(req.Method).To(Equal(http.MethodPost))
		err := req.ParseMultipartForm(32 << 20)
		Ω(err).ShouldNot(HaveOccurred())
		file, _, _ := req.FormFile("file")
		tr := tar.NewReader(file)

	loop:
		for {
			header, err := tr.Next()
			switch {
			case err == io.EOF:
				break loop
			case err != nil:
				Ω(err).ShouldNot(HaveOccurred())
			case header == nil:
				continue
			}
			filesFound = append(filesFound, header.Name)

		}
		defer req.Body.Close()
		rw.WriteHeader(http.StatusOK)
	}))
	projectsFolder := "../test_helpers/projects/publishProjects/"
	areEquals := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for _, value := range a {
			index := sort.SearchStrings(b, value)
			if b[index] != value {
				return false
			}
		}
		return true
	}
	moveToProjectDirectory := func(project string) {
		err := os.Chdir(projectsFolder + project)
		Ω(err).ShouldNot(HaveOccurred())
	}
	getPc := func() *config.ProjectConfig {
		projectRoot, err := filepath.Abs(".")
		Ω(err).ShouldNot(HaveOccurred())
		manifestPath := "./c3pm.yml"
		m, err := manifest.Load(manifestPath)
		Ω(err).ShouldNot(HaveOccurred())
		pc := &config.ProjectConfig{Manifest: m, ProjectRoot: projectRoot}
		return pc
	}

	BeforeEach(func() {
		var err error
		wd, err = os.Getwd()
		Ω(err).ShouldNot(HaveOccurred())
		filesFound = nil
	})
	AfterEach(func() {
		err := os.Chdir(wd)
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("basic publish without config: should have sources, includes and c3pm.yml", func() {
		moveToProjectDirectory("basic")
		expectedFiles := []string{"src/main.cpp", "c3pm.yml"}

		apiClient := api.New(server.Client(), "")

		err := os.Setenv("C3PM_API_ENDPOINT", server.URL)
		Ω(err).ShouldNot(HaveOccurred())
		defer os.Unsetenv("C3PM_API_ENDPOINT")

		err = ctpm.Publish(getPc(), apiClient)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(areEquals(expectedFiles, filesFound)).Should(BeTrue())
	})

	It("Include and Exclude files: should have toto.txt but not README.md in tarball", func() {
		moveToProjectDirectory("includeExclude")
		expectedFiles := []string{"toto.txt", "src/main.cpp", "c3pm.yml"}

		apiClient := api.New(server.Client(), "")

		err := os.Setenv("C3PM_API_ENDPOINT", server.URL)
		Ω(err).ShouldNot(HaveOccurred())
		defer os.Unsetenv("C3PM_API_ENDPOINT")

		err = ctpm.Publish(getPc(), apiClient)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(areEquals(expectedFiles, filesFound)).Should(BeTrue())
	})

	It("file in Include and Exclude: shouldn't have README.md in tarball", func() {
		moveToProjectDirectory("includeAndExcludeFile")
		expectedFiles := []string{"src/main.cpp", "c3pm.yml"}

		apiClient := api.New(server.Client(), "")

		err := os.Setenv("C3PM_API_ENDPOINT", server.URL)
		Ω(err).ShouldNot(HaveOccurred())
		defer os.Unsetenv("C3PM_API_ENDPOINT")

		err = ctpm.Publish(getPc(), apiClient)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(areEquals(expectedFiles, filesFound)).Should(BeTrue())
	})

	It("overwrite default include: shouldn't have c3pm.yml in tarball", func() {
		moveToProjectDirectory("excludeManifest")
		expectedFiles := []string{"src/main.cpp"}

		apiClient := api.New(server.Client(), "")

		err := os.Setenv("C3PM_API_ENDPOINT", server.URL)
		Ω(err).ShouldNot(HaveOccurred())
		defer os.Unsetenv("C3PM_API_ENDPOINT")

		err = ctpm.Publish(getPc(), apiClient)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(areEquals(expectedFiles, filesFound)).Should(BeTrue())
	})
})
