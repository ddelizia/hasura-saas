package oidc_test

import (
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz/oidc"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("", func() {

	logrus.SetOutput(ioutil.Discard)

	Context("GetAuthHeader()", func() {
		It("", func() {
			// Given
			req := &http.Request{
				Header: http.Header{
					"Authorization": []string{"Bearer some.bearer.token"},
				},
			}

			// When
			got := oidc.GetAuthHeader(req)

			// Then
			Expect(got).To(Equal("some.bearer.token"))
		})
	})

})
