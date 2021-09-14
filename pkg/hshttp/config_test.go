package hshttp_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [subscription]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigAccountIdHeader() not to be empty", func() {
		Expect(hshttp.ConfigAccountIdHeader()).ShouldNot(BeEmpty())
	})
})
