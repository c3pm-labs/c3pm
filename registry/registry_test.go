package registry_test

import (
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Add", func() {
	Describe("Registry", func() {
		var server *ghttp.Server
		var options registry.Options

		BeforeEach(func() {
			server = ghttp.NewServer()
			options.RegistryURL = server.URL()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/", "prefix=boost&typeList=2"),
					ghttp.RespondWithJSONEncoded(200, []string{"0.0.1", "1.0.0"}),
				),
			)
		})
		It("fetches the right version", func() {
			version, err := registry.GetLastVersion("boost", options)
			if err != nil {
				Fail(err.Error())
			}
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(version).Should(Equal(semver.MustParse("1.0.0")))
		})
	})
})
