package rp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/sirupsen/logrus"
)

type HasuraHandler struct {
	Proxy    http.Handler
	AuthzSvc authz.Service
}

func NewHasuraService(authzSvc authz.Service) http.Handler {
	return &HasuraHandler{
		Proxy:    httputil.NewSingleHostReverseProxy(ConfigEsUrl()),
		AuthzSvc: authzSvc,
	}
}

func (h *HasuraHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*2))
	defer cancel()
	req = req.WithContext(ctx)

	hshttp.SetSslRedirect(req, ConfigHasuraUrl())

	hshttp.SetHaderOnRequest(req, "X-Hasura-Admin-Secret", gqlreq.ConfigAdminSecret())

	req.Host = ConfigHasuraUrl().Host
	req.URL.Path = "/"

	authToken := req.Header.Get(hshttp.ConfigJWTHeader())

	accountId := req.Header.Get(hshttp.ConfigAccountIdHeader())

	authzInfo, err := h.AuthzSvc.GetAuthInfo(ctx, authToken, accountId)
	if err != nil {
		hshttp.WriteError(res, err)
		return

	} else {
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraAccountHeader(), accountId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraUserIdHeader(), authzInfo.UserId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraRoleHeader(), authzInfo.RoleId)

		logrus.Debug(fmt.Sprintf("user [%s], accessing tenant [%s], with role id [%s]", authzInfo.UserId, accountId, authzInfo.RoleId))
	}

	h.Proxy.ServeHTTP(res, req)

}
