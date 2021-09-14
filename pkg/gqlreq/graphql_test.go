package gqlreq_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test_Graphql(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "[graphql] test suite")
}
