package oidc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Authz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[oidc] test suite")
}
