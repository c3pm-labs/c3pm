package defaultadapter

import (
	"bufio"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("CMake interaction", func() {
	Describe("CMake build file generation", func() {
		const (
			BUILD_DIR = "/tmp/c3pm_cmake_test1"
		)

		AfterEach(func() {
			_ = os.RemoveAll(BUILD_DIR)
		})
		It("does generate the build directory", func() {
			err := cmakeGenerateBuildFiles("../../test_helpers/projects/cmakeProject", BUILD_DIR, map[string]string{})
			Ω(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(BUILD_DIR)
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("uses the variables added", func() {
			err := cmakeGenerateBuildFiles("../../test_helpers/projects/cmakeProject", BUILD_DIR, map[string]string{"CMAKE_AR:FILEPATH": "/bin/ls"})
			Ω(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(BUILD_DIR)
			Ω(err).ShouldNot(HaveOccurred())
			// To verify that the variable has been applied, check that it is contained in the CMakeCache.txt file in the build directory.
			// The file contains a line per variable, with the format NAME:TYPE=VALUE
			f, err := os.Open(filepath.Join(BUILD_DIR, "CMakeCache.txt"))
			Ω(err).ShouldNot(HaveOccurred())
			rd := bufio.NewReader(f)
			found := false
			var s string
			for s, err = "", nil; err == nil; s, err = rd.ReadString('\n') {
				if strings.Contains(s, "CMAKE_AR") && strings.Contains(s, "/bin/ls") {
					found = true
				}
			}
			Ω(found).Should(BeTrue())
		})
	})

	Describe("CMake cmakeBuild", func() {
		const (
			BUILD_DIR = "/tmp/c3pm_cmake_test2"
		)

		AfterEach(func() {
			_ = os.RemoveAll(BUILD_DIR)
		})
		It("builds the project", func() {
			// Generate files
			err := cmakeGenerateBuildFiles("../../test_helpers/projects/cmakeProject", BUILD_DIR, map[string]string{})
			Ω(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(BUILD_DIR)
			Ω(err).ShouldNot(HaveOccurred())

			// Build the project
			err = cmakeBuild(BUILD_DIR)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(filepath.Join(BUILD_DIR, "test1"))
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
