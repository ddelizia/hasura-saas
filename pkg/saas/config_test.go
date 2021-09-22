package saas_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/rp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [saas]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigListenAddress() not to be empty", func() {
		Expect(rp.ConfigListenAddress()).ShouldNot(BeEmpty())
	})

})
