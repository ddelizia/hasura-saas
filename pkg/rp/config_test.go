package rp_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/rp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [rp]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigEsUrl() not to be empty", func() {
		Expect(rp.ConfigEsUrl().Host).ShouldNot(BeEmpty())
	})

	It("ConfigEsPublicIndex() not to be empty", func() {
		Expect(rp.ConfigEsPublicIndex()).ShouldNot(BeEmpty())
	})

	It("ConfigEsPrivateIndex() not to be empty", func() {
		Expect(rp.ConfigEsPrivateIndex()).ShouldNot(BeEmpty())
	})

	It("ConfigEsAuthorizationHeader() not to be empty", func() {
		Expect(rp.ConfigEsAuthorizationHeader()).ShouldNot(BeEmpty())
	})

	It("ConfigHasuraUrl() not to be empty", func() {
		Expect(rp.ConfigHasuraUrl().Host).ShouldNot(BeEmpty())
	})

})
