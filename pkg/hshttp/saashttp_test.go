package hshttp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_hshttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[hshttp] test suite")
}
