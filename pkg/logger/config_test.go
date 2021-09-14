package logger_test

import (
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("config [logger]", func() {

	logrus.SetOutput(ioutil.Discard)

	It("ConfigLogLevel() not to be empty", func() {
		Expect(logger.ConfigLogLevel()).ShouldNot(BeEmpty())
	})

})
