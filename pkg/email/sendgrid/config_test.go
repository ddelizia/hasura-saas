package sendgrid_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/email/sendgrid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [sendgrid]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigSendgridApiKey() not to be empty", func() {
		Expect(sendgrid.ConfigSendgridApiKey()).ShouldNot(BeEmpty())
	})
})
