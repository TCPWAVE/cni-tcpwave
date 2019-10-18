package cniutils

import (
	"testing"
	gkg "github.com/onsi/ginkgo"
	gmg "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	gmg.RegisterFailHandler(gkg.Fail)
	gkg.RunSpecs(t, "Cni Utils Suite")
}
