package manifest

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
)

var exampleConfig = `
c3pm_version: v1
type: library
name: c3pm
description: This is the package c3pm
version: 1.0.0
license: ISC
files:
    sources:
        - 'src/**/*.cpp'
    includes:
        - include/private/header.h
        - include/private/private.h
    include_dirs:
        - include/private
    exports:
        - 'include/public/**/*.h'
dependencies:
    future: 12.2.3
    past: 2.0.0
`

func TestLoadManifest(t *testing.T) {
	g := NewGomegaWithT(t)
	pc, err := deserialize([]byte(exampleConfig))

	g.Expect(err).To(BeNil())
	g.Expect(pc.C3pmVersion).To(Equal(C3pmVersion1))
	g.Expect(pc.Type).To(Equal(Library))
	g.Expect(pc.Name).To(Equal("c3pm"))
	g.Expect(pc.Description).To(Equal("This is the package c3pm"))
	v, err := VersionFromString("1.0.0")
	g.Expect(err).To(BeNil())
	g.Expect(pc.Version).To(Equal(v))
	g.Expect(pc.License).To(Equal("ISC"))
	var m = make(map[string]string)
	m["future"] = "12.2.3"
	m["past"] = "2.0.0"
	g.Expect(pc.Dependencies["future"]).To(Equal(m["future"]))
	g.Expect(pc.Dependencies["past"]).To(Equal(m["past"]))
}

func TestSaveManifest(t *testing.T) {
	g := NewGomegaWithT(t)
	pc, err := deserialize([]byte(exampleConfig))

	g.Expect(err).To(BeNil())
	data, err := pc.serialize()
	g.Expect(err).To(BeNil())
	newPc, err := deserialize(data)
	g.Expect(err).To(BeNil(), "saved config was: %s", string(data))
	fmt.Println(newPc)
	fmt.Println(pc)
	g.Expect(newPc).To(Equal(pc))
}
