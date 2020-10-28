package manifest

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestC3pmVersionFromString(t *testing.T) {
	g := NewGomegaWithT(t)
	v, err := C3pmVersionFromString("v1")
	g.Expect(err).To(BeNil())
	g.Expect(v).To(Equal(C3pmVersion1))
}

func TestC3pmVersion_String(t *testing.T) {
	g := NewGomegaWithT(t)
	v, err := C3pmVersionFromString("v1")
	g.Expect(err).To(BeNil())
	g.Expect(v.String()).To(Equal("v1"))
}
