package config_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCtpm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}
