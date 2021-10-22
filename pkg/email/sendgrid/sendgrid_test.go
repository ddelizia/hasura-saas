package sendgrid_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[sendgrid] test suite")
}
