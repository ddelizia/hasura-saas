package hshttp_test

import (
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("test headers", func() {

	logrus.SetOutput(ioutil.Discard)

	Context("CleanHasuraSaasHeaders()", func() {
		It("should clean headers from the request", func() {
			// Given
			const hToDelete = "X-Hasura-Saas-Something"
			const hTokeep = "X-To-Keep"
			req := &http.Request{
				Header: http.Header{
					hToDelete: []string{"something"},
					hTokeep:   []string{"something else"},
				},
			}

			// When
			hshttp.CleanHasuraSaasHeaders(req)

			// Then
			Expect(req.Header.Get(hToDelete)).To(BeEmpty())
			Expect(req.Header.Get(hTokeep)).ToNot(BeEmpty())
		})
	})

})
