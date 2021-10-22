package email_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/email"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [email]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigFrom() not to be empty", func() {
		Expect(email.ConfigFrom()).ShouldNot(BeEmpty())
	})
})
