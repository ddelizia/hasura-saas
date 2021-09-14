package env_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Env(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[env] test suite")
}
