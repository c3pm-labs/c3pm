package config_test

import (
	"github.com/c3pm-labs/c3pm/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Authentication configuration", func() {
	var (
		TargetDir string
	)

	BeforeEach(func() {
		dir, err := ioutil.TempDir("", "c3pm_test_*")
		Ω(err).ShouldNot(HaveOccurred())
		TargetDir = dir
		err = ioutil.WriteFile(filepath.Join(TargetDir, "auth.cfg"), []byte("testtoken"), os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		_ = os.RemoveAll(TargetDir)
	})
	Context("TokenStrict", func() {
		It("gets the correct token", func() {
			err := os.Setenv("C3PM_USER_DIR", TargetDir)
			Ω(err).ShouldNot(HaveOccurred())
			token, err := config.TokenStrict()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).Should(Equal("testtoken"))
		})
	})
	Context("Token", func() {
		It("gets the correct token", func() {
			err := os.Setenv("C3PM_USER_DIR", TargetDir)
			Ω(err).ShouldNot(HaveOccurred())
			token := config.Token()
			Ω(token).Should(Equal("testtoken"))
		})
		It("returns an empty string in case of error", func() {
			err := os.Setenv("C3PM_USER_DIR", "/tmp/fakedir")
			Ω(err).ShouldNot(HaveOccurred())
			token := config.Token()
			Ω(token).Should(BeEmpty())
		})
	})
})
