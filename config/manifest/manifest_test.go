package manifest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var exampleConfig = `
c3pm_version: v1
type: library
name: c3pm
description: This is the package c3pm
version: 1.0.0
Documentation: "http://docs.c3pm.io/"
Website: "https://c3pm.io/"
Repository: "https://github.com/c3pm-labs"
Contributors: "Alex Hugh - Ramy J."
license: ISC
publish:
    include_dirs:
    - 'include/public'
build:
    adapter:
        name: c3pm
        version: 0.0.1
    config:
        sources:
        - 'src/**/*.cpp'
        headers:
        - include/private/header.h
        - include/private/private.h
        include_dirs:
        - include/private
dependencies:
    future: 12.2.3
    past: 2.0.0
test_dependencies:
    catch: 2.0.0
`

var _ = Describe("manifest", func() {
	It("test load of the manifest", func() {
		pc, err := deserialize([]byte(exampleConfig))
		Ω(err).ShouldNot(HaveOccurred())

		Ω(pc.C3PMVersion).To(Equal(C3PMVersion1))
		Ω(pc.Type).To(Equal(Library))
		Ω(pc.Name).To(Equal("c3pm"))
		Ω(pc.Description).To(Equal("This is the package c3pm"))
		v, err := VersionFromString("1.0.0")
		Ω(err).ShouldNot(HaveOccurred())

		Ω(pc.Version).To(Equal(v))
		Ω(pc.License).To(Equal("ISC"))
		var m = make(map[string]string)
		m["future"] = "12.2.3"
		m["past"] = "2.0.0"
		Ω(pc.Dependencies["future"]).To(Equal(m["future"]))
		Ω(pc.Dependencies["past"]).To(Equal(m["past"]))
	})

	It("Test save manifest", func() {
		pc, err := deserialize([]byte(exampleConfig))
		Ω(err).ShouldNot(HaveOccurred())

		data, err := pc.serialize()
		Ω(err).ShouldNot(HaveOccurred())

		newPc, err := deserialize(data)
		Ω(err).To(BeNil(), "saved config was: %s", string(data))
		Ω(newPc).To(Equal(pc))
	})
})
