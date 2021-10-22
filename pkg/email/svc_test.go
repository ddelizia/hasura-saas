package email_test

import (
	"errors"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/email"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

type EmailData struct {
	StringData string
	IntData    int64
	BoolData   bool
}

var _ = Describe("service", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		s email.Service
	)

	BeforeEach(func() {
		s = email.NewServiceMock()
	})

	mockSendEmail := func(e error) {
		s.(*email.ServiceMock).On("SendEmail", []string{"email1@something.com", "email2@something.com"},
			"This is the subject of the email",
			"This is the template id",
			&EmailData{
				StringData: "some string",
				IntData:    10,
				BoolData:   true,
			}).Times(1).Return(e)
	}

	It("should be able to setup multiple email sender implementations (mock) and send eamail correctly", func() {
		// Given
		mockSendEmail(nil)

		// When
		err := s.SendEmail(
			[]string{"email1@something.com", "email2@something.com"},
			"This is the subject of the email",
			"This is the template id",
			&EmailData{
				StringData: "some string",
				IntData:    10,
				BoolData:   true,
			},
		)

		// Then
		Expect(err).Should(BeNil())
	})

	It("should be able to setup multiple email sender implementations (mock) and should give error", func() {
		// Given
		mockSendEmail(errors.New("Error is thrown"))

		// When
		err := s.SendEmail(
			[]string{"email1@something.com", "email2@something.com"},
			"This is the subject of the email",
			"This is the template id",
			&EmailData{
				StringData: "some string",
				IntData:    10,
				BoolData:   true,
			},
		)

		// Then
		Expect(err).ShouldNot(BeNil())
	})

})
