package ctpm_test

import (
	"github.com/gabrielcolson/c3pm/cli/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path"
)

var _ = Describe("Logout", func() {
	c3pmHomeDir := getTestFolder("LogoutUserHome")

	It("should delete auth file", func() {
		err := os.MkdirAll(c3pmHomeDir, os.ModePerm)
		Ω(err).To(BeNil())
		err = os.Setenv("C3PM_USER_DIR", c3pmHomeDir)
		Ω(err).To(BeNil())
		defer os.Unsetenv("C3PM_USER_DIR")

		authFilePath := path.Join(c3pmHomeDir, "auth.cfg")
		f, err := os.Create(authFilePath)
		Ω(err).To(BeNil())
		err = f.Close()
		Ω(err).To(BeNil())

		err = ctpm.Logout()
		Ω(err).To(BeNil())

		_, err = os.Stat(authFilePath)
		Ω(os.IsNotExist(err)).To(BeTrue())
	})
})
