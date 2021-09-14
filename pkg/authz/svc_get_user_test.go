package authz_test

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("getUserImpl.GetUser()", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		s *authz.UserGetterImpl
	)

	BeforeEach(func() {
		s = &authz.UserGetterImpl{}
	})

	Context("Anything", func() {

		It("should return the user when validation is successful", func() {
			// Given
			authz.ValidateJwtFunc = func(jwtB64 string) (*authz.JwtTokenWithClaims, error) {
				return &authz.JwtTokenWithClaims{StandardClaims: &jwt.StandardClaims{Subject: "user"}}, nil
			}

			// When
			got, err := s.GetUser(context.Background(), "some.jwt.token")

			// Then
			Expect(got).To(Equal("user"))
			Expect(err).To(BeNil())
		})

		It("should return an error when validation of the token has failed", func() {
			// Given
			authz.ValidateJwtFunc = func(jwtB64 string) (*authz.JwtTokenWithClaims, error) {
				return nil, errors.New("some validation error")
			}

			// When
			got, err := s.GetUser(context.Background(), "some.jwt.token")

			// Then
			Expect(got).To(Equal(""))
			Expect(err).NotTo(BeNil())
		})

		It("should return the anonymous user when there is no token", func() {
			// When
			got, err := s.GetUser(context.Background(), "")

			// Then
			Expect(got).To(Equal(authz.ConfigAnonymousUser()))
			Expect(err).To(BeNil())
		})
	})

})
