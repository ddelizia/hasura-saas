package oidc_test

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/authz/oidc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("AccountGetterOidc", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		s *oidc.AccountGetterOidc
	)

	BeforeEach(func() {
		s = &oidc.AccountGetterOidc{}
	})

	Context("GetAccount()", func() {

		It("should return the account when validation is successful", func() {
			// Given
			httpReq := &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer somejwt"},
				},
			}
			oidc.ValidateJwtFunc = func(jwtB64 string, r *http.Request) (*oidc.JwtTokenWithClaims, error) {
				return &oidc.JwtTokenWithClaims{HasuraSaasNamespace: oidc.HasuraSaasClaims{HasuraSaasAccount: "account"}}, nil
			}

			// When
			got, err := s.GetAccount(httpReq)

			// Then
			Expect(got).To(Equal("account"))
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
			got, err := s.GetAccount(httpReq)

			// Then
			Expect(got).To(Equal(""))
			Expect(err).NotTo(BeNil())
		})

		It("should return the anonymous account when there is no token", func() {
			// When
			httpReq := &http.Request{
				Header: http.Header{},
			}
			got, err := s.GetAccount(httpReq)

			// Then
			Expect(got).To(Equal(authz.ConfigAnonymousRole()))
			Expect(err).To(BeNil())
		})
	})

})
