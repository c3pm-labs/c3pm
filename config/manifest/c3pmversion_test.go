package manifest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("c3pm version Test", func() {
	It("Check C3PM version from string", func() {
		v, err := C3PMVersionFromString("v1")
		Ω(err).ShouldNot(HaveOccurred())

		Ω(v).To(Equal(C3PMVersion1))
	})
})
