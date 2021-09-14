package authz_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/hstest"
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
		roleGetterMock *authz.RoleGetterMock
		userGetterMock *authz.UserGetterMock
		s              *authz.AuthInfoGetterImpl
	)

	mockGetRole := func(role string, err error, arguments ...interface{}) {
		roleGetterMock.On("GetRole", arguments...).Return(role, err)
	}

	mockGetUser := func(user string, err error, arguments ...interface{}) {
		userGetterMock.On("GetUser", arguments...).Return(user, err)
	}

	BeforeEach(func() {
		roleGetterMock = &authz.RoleGetterMock{}
		userGetterMock = &authz.UserGetterMock{}
		s = &authz.AuthInfoGetterImpl{
			UserGetter: userGetterMock,
			RoleGetter: roleGetterMock,
		}

	})

	Context("GetAuthInfo()", func() {
		It("should return data from role and user service", func() {
			// Given
			mockGetRole("role", nil, mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("user", nil, mock.Anything, mock.Anything)

			// When
			got, err := s.GetAuthInfo(context.Background(), JWT, "account")

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(&authz.AuthzInfo{
				RoleId:    "role",
				UserId:    "user",
				AccountId: "account",
				Jwt:       JWT,
			}))
		})

		It("should throw an error when GetRole call returns an error", func() {
			// Given
			mockGetRole("", errors.New("some error"), mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("user", nil, mock.Anything, mock.Anything)

			// When
			_, err := s.GetAuthInfo(context.Background(), JWT, "account")

			// Then
			Expect(err).To(Not(BeNil()))
		})

		It("should throw an error when GetUser call returns an error", func() {
			// Given
			mockGetRole("role", nil, mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("", errors.New("some error"), mock.Anything, mock.Anything)

			// When
			_, err := s.GetAuthInfo(context.Background(), JWT, "account")

			// Then
			Expect(err).To(Not(BeNil()))
		})
	})

	Context("GetAuthInfoFromRequest()", func() {

		var req *http.Request

		BeforeEach(func() {
			req = hstest.CrerateRequest("POST", "/some/path", nil, nil)
			req.Header.Set(hshttp.AccountHeaderName(), "account")
			req.Header.Set(hshttp.JwtHeaderName(), JWT)
		})

		It("should return data from role and user service", func() {
			// Given
			mockGetRole("role", nil, mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("user", nil, mock.Anything, mock.Anything)

			// When
			got, err := s.GetAuthInfoFromRequest(req)

			// Then
			Expect(err).To(BeNil())
			Expect(got).To(Equal(&authz.AuthzInfo{
				RoleId:    "role",
				UserId:    "user",
				AccountId: "account",
				Jwt:       JWT,
			}))
		})

		It("should throw an error when GetRole call returns an error", func() {
			// Given
			mockGetRole("", errors.New("some error"), mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("user", nil, mock.Anything, mock.Anything)

			// When
			_, err := s.GetAuthInfoFromRequest(req)

			// Then
			Expect(err).To(Not(BeNil()))
		})

		It("should throw an error when GetUser call returns an error", func() {
			// Given
			mockGetRole("role", nil, mock.Anything, mock.Anything, mock.Anything)
			mockGetUser("", errors.New("some error"), mock.Anything, mock.Anything)

			// When
			_, err := s.GetAuthInfoFromRequest(req)

			// Then
			Expect(err).To(Not(BeNil()))
		})
	})
})
