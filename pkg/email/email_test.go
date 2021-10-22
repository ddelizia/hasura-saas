package email_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[email] test suite")
}
