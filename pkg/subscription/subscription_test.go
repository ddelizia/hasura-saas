package subscription_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Subscription(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[subscription] test suite")
}
