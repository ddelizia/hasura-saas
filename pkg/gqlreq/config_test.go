package gqlreq_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [graphql]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigGraphqlURL() not to be empty", func() {
		Expect(gqlreq.ConfigGraphqlURL()).ShouldNot(BeEmpty())
	})

	It("ConfigAdminSecret() not to be empty", func() {
		Expect(gqlreq.ConfigAdminSecret()).ShouldNot(BeEmpty())
	})

	It("ConfigHasuraUserIdHeader() not to be empty", func() {
		Expect(gqlreq.ConfigHasuraUserIdHeader()).ShouldNot(BeEmpty())
	})

	It("ConfigHasuraRoleHeader() not to be empty", func() {
		Expect(gqlreq.ConfigHasuraRoleHeader()).ShouldNot(BeEmpty())
	})

	It("ConfigHasuraAccountHeader() not to be empty", func() {
		Expect(gqlreq.ConfigHasuraAccountHeader()).ShouldNot(BeEmpty())
	})

})
