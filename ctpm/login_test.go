package ctpm_test

import (
	"encoding/json"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/ctpm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
)

var _ = Describe("Login", func() {
	c3pmHomeDir := getTestFolder("LoginUserHome")
	testLogin := "demo@demo.com"
	testPassword := "demodemo"
	testApiKey := "demo"

	It("should login without error", func() {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			Ω(req.URL.String()).To(Equal("/auth/login"))
			Ω(req.Method).To(Equal(http.MethodPost))

			type loginPayload struct {
				Login    string `json:"login"`
				Password string `json:"password"`
			}
			want := loginPayload{testLogin, testPassword}

			var got loginPayload
			body, err := ioutil.ReadAll(req.Body)
			Ω(err).To(BeNil())
			err = json.Unmarshal(body, &got)
			Ω(err).To(BeNil())

			Ω(want).To(Equal(got))

			res, err := json.Marshal(struct{ ApiKey string }{testApiKey})
			Ω(err).To(BeNil())
			_, err = rw.Write(res)
			Ω(err).To(BeNil())
		}))
		defer server.Close()
		apiClient := api.API{Client: server.Client()}

		err := os.MkdirAll(c3pmHomeDir, os.ModePerm)
		Ω(err).To(BeNil())
		err = os.Setenv("C3PM_USER_DIR", c3pmHomeDir)
		Ω(err).To(BeNil())
		defer os.Unsetenv("C3PM_USER_DIR")

		err = os.Setenv("C3PM_API_ENDPOINT", server.URL)
		Ω(err).To(BeNil())
		defer os.Unsetenv("C3PM_API_ENDPOINT")

		err = ctpm.Login(apiClient, testLogin, testPassword)
		Ω(err).To(BeNil())
	})

	It("should create auth file", func() {
		got, err := ioutil.ReadFile(path.Join(c3pmHomeDir, "auth.cfg"))
		Ω(err).To(BeNil())
		Ω(string(got)).To(Equal(testApiKey))
	})
})
