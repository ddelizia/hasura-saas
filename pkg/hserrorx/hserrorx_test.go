package hserrorx_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_hserrorx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[hserrorx] test suite")
}
