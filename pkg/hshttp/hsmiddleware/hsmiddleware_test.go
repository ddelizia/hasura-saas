package hsmiddleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Hsmiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[hsmiddleware] test suite")
}