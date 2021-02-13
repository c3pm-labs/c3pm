package ctpm_test

import (
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path"
)

var _ = Describe("Logout", func() {
	c3pmHomeDir := getTestFolder("LogoutUserHome")

	It("should delete auth file", func() {
		err := os.MkdirAll(c3pmHomeDir, os.ModePerm)
		Ω(err).ShouldNot(HaveOccurred())

		err = os.Setenv("C3PM_USER_DIR", c3pmHomeDir)
		Ω(err).ShouldNot(HaveOccurred())

		defer os.Unsetenv("C3PM_USER_DIR")

		authFilePath := path.Join(c3pmHomeDir, "auth.cfg")
		f, err := os.Create(authFilePath)
		Ω(err).ShouldNot(HaveOccurred())

		err = f.Close()
		Ω(err).ShouldNot(HaveOccurred())

		err = ctpm.Logout()
		Ω(err).ShouldNot(HaveOccurred())

		_, err = os.Stat(authFilePath)
		Ω(os.IsNotExist(err)).To(BeTrue())
	})
})
