package manifest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dependencies", func() {
	It("test dependencies map", func() {
		var m = make(map[string]string)
		m["boost"] = "1.2.3"
		d, err := DependenciesFromMap(m)
		Ω(err).To(BeNil())
		Ω(d).To(Equal(Dependencies(m)))
	})
})
