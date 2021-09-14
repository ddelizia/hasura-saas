package subscription_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/subscription"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [subscription]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigWebhookSecret() not to be empty", func() {
		Expect(subscription.ConfigWebhookSecret()).ShouldNot(BeEmpty())
	})

	It("ConfigDomain() not to be empty", func() {
		Expect(subscription.ConfigDomain()).ShouldNot(BeEmpty())
	})

	It("ConfigWebhookListenAddress() not to be empty", func() {
		Expect(subscription.ConfigWebhookListenAddress()).ShouldNot(BeEmpty())
	})

	It("ConfigApiKey() not to be empty", func() {
		Expect(subscription.ConfigApiKey()).ShouldNot(BeEmpty())
	})

})
