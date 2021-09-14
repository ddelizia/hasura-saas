package authz_test

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("getRoleImpl.GetRole()", func() {

	logrus.SetOutput(ioutil.Discard)

	const (
		ADMIN_ROLE = "admin"
	)

	var (
		graphqlMock *gqlreq.ServiceMock
		s           *authz.RoleGetterImpl
	)

	BeforeEach(func() {
		graphqlMock = gqlreq.NewServiceMock().(*gqlreq.ServiceMock)
		s = &authz.RoleGetterImpl{
			GraphQlSvc: graphqlMock,
		}

	})

	mockGraphqlExecute := func(err error, runFn func(args mock.Arguments), arguments ...interface{}) {
		graphqlMock.ExecuterMock.On("Execute", arguments...).Run(runFn).Return(err)
	}

	Context("Anything", func() {
		It("should get role returned by graphQL", func() {
			// Given
			mockGraphqlExecute(nil, func(args mock.Arguments) {
				data := args.Get(5).(*authz.SaasMembershipResponse)
				data.SaasMembership = []authz.SaasMembership{
					{SaasRole: authz.SaasRole{HasuraRole: ADMIN_ROLE}},
				}
			}, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

			// When
			got, err := s.GetRole(context.Background(), "account", "user")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(ADMIN_ROLE))
		})

		It("should get role user role when user is present and cannot find the role in hasura", func() {
			// Given
			mockGraphqlExecute(nil, func(args mock.Arguments) {
				data := args.Get(5).(*authz.SaasMembershipResponse)
				data.SaasMembership = []authz.SaasMembership{}
			}, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

			// When
			got, err := s.GetRole(context.Background(), "account", "user")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(authz.ConfigLoggedInRole()))
		})

		It("should get role anonymous when user is not specified", func() {
			// When
			got, err := s.GetRole(context.Background(), "account", "")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(authz.ConfigAnonymousRole()))
		})

		It("should get role user when account is not specified", func() {
			// When
			got, err := s.GetRole(context.Background(), "", "user")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(authz.ConfigLoggedInRole()))
		})

		It("should get error when multiple value are returned from hasura", func() {
			// Given
			mockGraphqlExecute(errors.New("some error"), func(args mock.Arguments) {
				data := args.Get(5).(*authz.SaasMembershipResponse)
				data.SaasMembership = []authz.SaasMembership{
					{SaasRole: authz.SaasRole{HasuraRole: "admin"}},
					{SaasRole: authz.SaasRole{HasuraRole: "user"}},
				}
			}, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)

			// When
			_, err := s.GetRole(context.Background(), "account", "user")

			// Then
			Expect(err).To(Not(BeNil()))
		})

	})
})
