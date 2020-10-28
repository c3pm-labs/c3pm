package manifest

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestDependencies_Map(t *testing.T) {
	g := NewGomegaWithT(t)
	var m = make(map[string]string)
	m["boost"] = "1.2.3"
	d, err := DependenciesFromMap(m)
	g.Expect(err).To(BeNil())
	g.Expect(d).To(Equal(Dependencies(m)))
}
