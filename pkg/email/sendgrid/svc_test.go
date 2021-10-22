package sendgrid_test

import (
	"io/ioutil"
	"os"

	"github.com/ddelizia/hasura-saas/pkg/email/sendgrid"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Email Sender", func() {

	logrus.SetOutput(ioutil.Discard)

	It("should send an email correctly in sandbox mode", func() {
		// Given
		os.Setenv("EMAIL.FROM", os.Getenv("SENDGRID_FROM"))
		os.Setenv("EMAIL.SENDGRID.APIKEY", os.Getenv("SENDGRID_API_KEY"))
		s := sendgrid.NewServiceSandbox()

		// When
		err := s.SendEmail(
			[]string{"email@sandbox.com"},
			"This is an example email",
			"d-6cdd181da221464a827c3170a363f427",
			map[string]string{
				"example": "example",
			})

		// Then
		Expect(err).Should(BeNil())
	})

})
