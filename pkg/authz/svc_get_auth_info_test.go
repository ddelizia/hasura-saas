package authz_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("getAuthInfoImpl", func() {
	logrus.SetOutput(ioutil.Discard)

	const (
		JWT = "some.jwt.token"
	)

	var (
		roleGetterMock    *authz.RoleGetterMock
		userGetterMock    *authz.UserGetterMock
		accountGetterMock *authz.AccountGetterMock
		s                 *authz.AuthInfoGetterImpl
	)

	mockGetRole := func(role string, err error, arguments ...interface{}) {
		roleGetterMock.On("GetRole", arguments...).Return(role, err)
	}

	mockGetUser := func(user string, err error, arguments ...interface{}) {
		userGetterMock.On("GetUser", arguments...).Return(user, err)
	}

	mockGetAccount := func(account string, err error, arguments ...interface{}) {
		accountGetterMock.On("GetAccount", arguments...).Return(account, err)
	}

	BeforeEach(func() {
		roleGetterMock = &authz.RoleGetterMock{}
		userGetterMock = &authz.UserGetterMock{}
		accountGetterMock = &authz.AccountGetterMock{}
		s = &authz.AuthInfoGetterImpl{
			UserGetter:    userGetterMock,
			RoleGetter:    roleGetterMock,
			AccountGetter: accountGetterMock,
		}

	})

	Context("GetAuthInfo()", func() {
		It("should return data from role and user service", func() {
			// Given
			mockGetUser("user", nil, mock.Anything)
			mockGetRole("role", nil, mock.Anything)
			mockGetAccount("account", nil, mock.Anything)

			// When
			got, err := s.GetAuthInfo(&http.Request{})

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(&authz.AuthzInfo{
				RoleId:    "role",
				UserId:    "user",
				AccountId: "account",
			}))
		})

		It("should throw an error when GetRole call returns an error", func() {
			// Given
			mockGetUser("user", nil, mock.Anything)
			mockGetRole("", errors.New("some error"), mock.Anything)
			mockGetAccount("account", nil, mock.Anything)

			// When
			_, err := s.GetAuthInfo(&http.Request{})

			// Then
			Expect(err).To(Not(BeNil()))
		})

		It("should throw an error when GetUser call returns an error", func() {
			// Given
			mockGetUser("", errors.New("some error"), mock.Anything)
			mockGetRole("role", nil, mock.Anything)
			mockGetAccount("account", nil, mock.Anything)

			// When
			_, err := s.GetAuthInfo(&http.Request{})

			// Then
			Expect(err).To(Not(BeNil()))
		})

		It("should throw an error when GetAccount call returns an error", func() {
			// Given
			mockGetUser("user", nil, mock.Anything)
			mockGetRole("role", nil, mock.Anything)
			mockGetAccount("", errors.New("some error"), mock.Anything)

			// When
			_, err := s.GetAuthInfo(&http.Request{})

			// Then
			Expect(err).To(Not(BeNil()))
		})
	})
})
