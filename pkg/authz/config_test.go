package authz_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [authz]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigJwksUrl() not to be empty", func() {
		Expect(authz.ConfigJwksUrl()).ShouldNot(BeEmpty())
	})

	It("ConfigAnonymousUser() not to be empty", func() {
		Expect(authz.ConfigAnonymousUser()).ShouldNot(BeEmpty())
	})

	It("ConfigAnonymousRole() not to be empty", func() {
		Expect(authz.ConfigAnonymousRole()).ShouldNot(BeEmpty())
	})

	It("ConfigLoggedInRole() not to be empty", func() {
		Expect(authz.ConfigLoggedInRole()).ShouldNot(BeEmpty())
	})

	It("ConfigAdminRole() not to be empty", func() {
		Expect(authz.ConfigAdminRole()).ShouldNot(BeEmpty())
	})

	It("ConfigAccountOwnerRole() not to be empty", func() {
		Expect(authz.ConfigAccountOwnerRole()).ShouldNot(BeEmpty())
	})

})
