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

var _ = Describe("handler_hasura.ServeHTTP()", func() {

	logrus.SetOutput(ioutil.Discard)

	var (
		authzMock *authz.ServiceMock
		proxyMock *rp.ReverseProxyMock
		handler   *rp.HasuraHandler
		rr        *httptest.ResponseRecorder
		req       *http.Request
	)

	BeforeEach(func() {
		authzMock = authz.NewServiceMock().(*authz.ServiceMock)
		proxyMock = rp.NewReverseProxyMock().(*rp.ReverseProxyMock)
		handler = &rp.HasuraHandler{
			Proxy:    proxyMock,
			AuthzSvc: authzMock,
		}
		rr = httptest.NewRecorder()
		req = hstest.CrerateRequest("POST", "/graphql/", nil, nil)
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
		It("should set hasura user and role when authz returns", func() {
			// Given
			req.Header.Set(hshttp.ConfigAccountIdHeader(), "account")
			mockGetAuthInfo(&authz.AuthzInfo{
				RoleId:    "admin",
				UserId:    "user",
				AccountId: "account",
			}, nil, mock.Anything, mock.Anything, mock.Anything)
			mockServeHTTP(mock.Anything, mock.Anything)

			// When
			handler.ServeHTTP(rr, req)

			// Then
			proxyRequest := proxyMock.Calls[0].Arguments[1].(*http.Request)
			Expect(proxyRequest.Header.Get("X-Hasura-User-Id")).Should(Equal("user"))
			Expect(proxyRequest.Header.Get("X-Hasura-Account-Id")).Should(Equal("account"))
			Expect(proxyRequest.Header.Get("X-Hasura-Role")).Should(Equal("admin"))
		})

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
