package api_test

import (
	"github.com/c3pm-labs/c3pm/api"
	apitest "github.com/c3pm-labs/c3pm/api/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
)

var _ = Describe("Upload", func() {
	var client api.API

	Context("server returning an error", func() {
		var srv *httptest.Server

		BeforeEach(func() {
			srv = httptest.NewServer(apitest.ErrorServer())
			client = api.New(srv.Client(), "xxx")
		})
		AfterEach(func() {
			srv.Close()
		})

		It("doesn't have enough storage available", func() {
			err := client.Upload([]string{"../test_helpers/main.cpp"})
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("server working correctly", func() {
		var srv *httptest.Server

		BeforeEach(func() {
			srv = httptest.NewServer(apitest.MockServer())
			client = api.New(srv.Client(), "xxx")
		})
		AfterEach(func() {
			srv.Close()
		})

		It("uploads the file correctly", func() {
			err := client.Upload([]string{"../test_helpers/main.cpp"})
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
