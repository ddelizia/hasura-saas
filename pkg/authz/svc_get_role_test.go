package authz_test

import (
	"context"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
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
		graphqlMock *gqlsdk.ServiceMock
		s           *authz.RoleGetterImpl
	)

	BeforeEach(func() {
		graphqlMock = gqlsdk.NewServiceMock().(*gqlsdk.ServiceMock)
		s = &authz.RoleGetterImpl{
			GraphQlSvc: graphqlMock,
		}
	})

	Context("Anything", func() {
		It("should get role returned by graphQL", func() {
			// Given
			graphqlMock.On("GetRoleForUserAndAccount", mock.Anything, "user", "account", mock.Anything).Return(
				&gqlsdk.QueryGetRoleForUserAndAccount{
					SaasMembership: []*struct {
						IDRole string "json:\"id_role\" graphql:\"id_role\""
					}{
						{
							IDRole: ADMIN_ROLE,
						},
					},
				},
				nil,
			)

			// When
			got, err := s.GetRole(context.Background(), "account", "user")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(ADMIN_ROLE))
		})

		It("should get role user role when user is present and cannot find the role in hasura", func() {
			// Given
			graphqlMock.On("GetRoleForUserAndAccount", mock.Anything, "user", "account", mock.Anything).Return(
				&gqlsdk.QueryGetRoleForUserAndAccount{
					SaasMembership: []*struct {
						IDRole string "json:\"id_role\" graphql:\"id_role\""
					}{},
				},
				nil,
			)

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
			graphqlMock.On("GetRoleForUserAndAccount", mock.Anything, "user", "account", mock.Anything).Return(
				&gqlsdk.QueryGetRoleForUserAndAccount{
					SaasMembership: []*struct {
						IDRole string "json:\"id_role\" graphql:\"id_role\""
					}{
						{
							IDRole: "admin",
						},
						{
							IDRole: "user",
						},
					},
				},
				nil,
			)

			// When
			_, err := s.GetRole(context.Background(), "account", "user")

			// Then
			Expect(err).To(Not(BeNil()))
		})

	})
})
