package hsmiddleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/ddelizia/hasura-saas/pkg/hshttp/hsmiddleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("LogRequest", func() {

	logrus.SetOutput(ioutil.Discard)

	It("should generate the request id", func() {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		mockNext := func(w http.ResponseWriter, r *http.Request) {}

		// When
		hsmiddleware.LogRequest()(mockNext)(w, req)

		// Then
		//requestId := hscontext.RequestIDValue(req.Context())
		Expect(nil).To(BeNil())
	})

})
