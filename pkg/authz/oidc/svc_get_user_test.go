package oidc_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/authz/oidc"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("UserGetterOidc", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		s *oidc.UserGetterOidc
	)

	BeforeEach(func() {
		s = &oidc.UserGetterOidc{}
	})

	Context("GetUser()", func() {

		It("should return the user when validation is successful", func() {
			// Given
			httpReq := &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer somejwt"},
				},
			}
			oidc.ValidateJwtFunc = func(jwtB64 string, r *http.Request) (*oidc.JwtTokenWithClaims, error) {
				return &oidc.JwtTokenWithClaims{StandardClaims: &jwt.StandardClaims{Subject: "user"}}, nil
			}

			// When
			got, err := s.GetUser(httpReq)

			// Then
			Expect(got).To(Equal("user"))
			Expect(err).To(BeNil())
		})

		It("should return an error when validation of the token has failed", func() {
			// Given
			httpReq := &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer some.jwt.token"},
				},
			}
			oidc.ValidateJwtFunc = func(jwtB64 string, r *http.Request) (*oidc.JwtTokenWithClaims, error) {
				return nil, errors.New("some validation error")
			}

			// When
			got, err := s.GetUser(httpReq)

			// Then
			Expect(got).To(Equal(""))
			Expect(err).NotTo(BeNil())
		})

		It("should return the anonymous user when there is no token", func() {
			// When
			httpReq := &http.Request{
				Header: http.Header{},
			}
			got, err := s.GetUser(httpReq)

			// Then
			Expect(got).To(Equal(authz.ConfigAnonymousUser()))
			Expect(err).To(BeNil())
		})
	})

})
