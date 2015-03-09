package ciface_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCiface(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ciface Suite")
}
