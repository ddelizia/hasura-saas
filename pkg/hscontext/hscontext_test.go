package hscontext_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Hscontext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[hscontext] test suite")
}
