package env_test

import (
	"io/ioutil"
	"os"

	"github.com/ddelizia/hasura-saas/pkg/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("index", func() {
	logrus.SetOutput(ioutil.Discard)

	Context("GetString()", func() {

		It("should get config from config.yaml", func() {
			Expect(env.GetString("env.mock.data")).Should(Equal("some value"))
		})

		It("shold load configuration from environemnt variable", func() {
			os.Setenv("ENV.MOCK.DATA", "some other value")
			Expect(env.GetString("env.mock.data")).Should(Equal("some other value"))
		})

	})
})
