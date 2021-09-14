package rp_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/hstest"
	"github.com/ddelizia/hasura-saas/pkg/rp"
	"github.com/joomcode/errorx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("handler_es.ServeHTTP()", func() {
	logrus.SetOutput(ioutil.Discard)

	OPERATION := "_msearch"

	var (
		authzMock *authz.ServiceMock
		proxyMock *rp.ReverseProxyMock
		handler   *rp.EsHandler
		rr        *httptest.ResponseRecorder
		req       *http.Request
	)

	BeforeEach(func() {
		authzMock = authz.NewServiceMock().(*authz.ServiceMock)
		proxyMock = rp.NewReverseProxyMock().(*rp.ReverseProxyMock)
		handler = &rp.EsHandler{
			Proxy:    proxyMock,
			AuthzSvc: authzMock,
		}
		rr = httptest.NewRecorder()
		req = hstest.CrerateRequest("POST", "/es/index/"+OPERATION, nil, map[string]string{"operation": OPERATION})
	})

	AfterEach(func() {

	})

	mockGetAuthInfo := func(authInfo *authz.AuthzInfo, err error, arguments ...interface{}) {
		authzMock.AuthInfoGetterMock.On("GetAuthInfo", arguments...).Return(authInfo, err)
	}

	mockServeHTTP := func(arguments ...interface{}) {
		proxyMock.On("ServeHTTP", arguments...).Return(nil)
	}

	Context("Any request is coming in", func() {
		It("should set the X-Forwarded-Host on the request", func() {
			// Given
			expectedHost := "this is the host"
			req.Header.Set("Host", expectedHost)

			mockGetAuthInfo(&authz.AuthzInfo{RoleId: authz.ConfigAnonymousRole()}, nil, mock.Anything, mock.Anything, mock.Anything)
			mockServeHTTP(mock.Anything, mock.Anything)

			// When
			handler.ServeHTTP(rr, req)

			// Then
			Expect(req.Header.Get("X-Forwarded-Host")).Should(Equal(expectedHost))
		})
	})

	Context("An HTTP request has user and account", func() {
		It("should redirect the public index", func() {
			// Given
			mockGetAuthInfo(&authz.AuthzInfo{RoleId: authz.ConfigAnonymousRole()}, nil, mock.Anything, mock.Anything, mock.Anything)
			mockServeHTTP(mock.Anything, mock.Anything)

			// When
			handler.ServeHTTP(rr, req)

			// Then
			proxyRequest := proxyMock.Calls[0].Arguments[1].(*http.Request)
			Expect(proxyRequest.URL.Path).Should(Equal(rp.ConfigEsPublicIndex() + OPERATION))
		})
	})

	Context("An HTTP request has no user account", func() {
		It("should redirect to the private index", func() {
			// Given
			mockGetAuthInfo(&authz.AuthzInfo{RoleId: "role", UserId: "userId"}, nil, mock.Anything, mock.Anything, mock.Anything)
			mockServeHTTP(mock.Anything, mock.Anything)

			// When
			handler.ServeHTTP(rr, req)

			// Then
			proxyRequest := proxyMock.Calls[0].Arguments[1].(*http.Request)
			Expect(proxyRequest.URL.Path).Should(Equal(rp.ConfigEsPrivateIndex() + OPERATION))
		})
	})

	Context("An Error occurs", func() {
		It("should get a response error when error is thrown by authz service", func() {
			// Given
			req.Header.Set(hshttp.ConfigAccountIdHeader(), "account")
			mockGetAuthInfo(nil, errorx.TimeoutElapsed.New("some error"), mock.Anything, mock.Anything, mock.Anything)
			mockServeHTTP(mock.Anything, mock.Anything)

			// When
			handler.ServeHTTP(rr, req)

			// Then
			res := rr.Result()
			defer res.Body.Close()

			data, _ := ioutil.ReadAll(res.Body)
			Expect(res.StatusCode).Should(Equal(http.StatusRequestTimeout))
			Expect(string(data)).Should(Equal("{\"error\":\"common.timeout: some error\"}\n"))
		})
	})

})
