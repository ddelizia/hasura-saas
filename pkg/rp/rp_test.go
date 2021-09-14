package rp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Rp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[rp] test suite")
}
