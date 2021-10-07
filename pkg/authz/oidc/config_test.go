package oidc_test

import (
	"io/ioutil"

	oidc "github.com/ddelizia/hasura-saas/pkg/authz/oidc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [authz]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigJwksUrl() not to be empty", func() {
		Expect(oidc.ConfigJwksUrl()).ShouldNot(BeEmpty())
	})

	It("ConfigJwksUrl() not to be empty", func() {
		Expect(oidc.ConfigJwksUrl()).ShouldNot(BeEmpty())
	})

})
