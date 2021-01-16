package registry_test

import (
	"encoding/xml"
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
			versions := registry.ListRegistryResponse{
				Name: "versions",
				Contents: []struct {
					Key string `xml:"Key"`
				}{{Key: "0.0.1"}, {Key: "1.0.0"}},
			}
			data, err := xml.Marshal(versions)
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/", "prefix=boost&typeList=2"),
					ghttp.RespondWith(200, data),
				),
			)
		})
		It("fetches the right version", func() {
			version, err := registry.GetLastVersion("boost", options)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(version).Should(Equal(semver.MustParse("1.0.0")))
		})
	})
})
