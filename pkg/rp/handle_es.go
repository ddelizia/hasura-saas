package rp

import (
	"context"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type EsHandler struct {
	Proxy    http.Handler
	AuthzSvc authz.Service
}

func NewEsHandler(authzSvc authz.Service) http.Handler {
	return &EsHandler{
		Proxy:    httputil.NewSingleHostReverseProxy(ConfigEsUrl()),
		AuthzSvc: authzSvc,
	}
}

func (h *EsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	esOperation := vars["operation"]
	logrus.Debug("operation requested [", esOperation, "]")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	req = req.WithContext(ctx)

	hshttp.SetSslRedirect(req, ConfigEsUrl())

	req.URL.Path = ConfigEsPublicIndex() + esOperation

	authToken := req.Header.Get(ConfigEsAuthorizationHeader())

	logrus.Debug("found token ", authToken)
	accountId := req.Header.Get(hshttp.ConfigAccountIdHeader())

	authoInfo, err := h.AuthzSvc.GetAuthInfo(ctx, authToken, accountId)
	if err != nil {
		logrus.Error("authorization info is not valid")
		hshttp.WriteError(res, err)
		return
	}

	if authoInfo.RoleId != authz.ConfigAnonymousRole() {
		req.URL.Path = ConfigEsPrivateIndex() + esOperation
	}

	h.Proxy.ServeHTTP(res, req)

}
