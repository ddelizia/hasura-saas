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

	hshttp.CleanHasuraSaasHeaders(req)

	hshttp.SetSslRedirect(req, ConfigHasuraUrl())

	req.Host = ConfigHasuraUrl().Host
	req.URL.Path = "/"

	accountId := req.Header.Get(hshttp.ConfigAccountIdHeader())

	authzInfo, err := h.AuthzSvc.GetAuthInfo(req)
	if err != nil {
		hshttp.WriteError(res, err)
		return

	} else {
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraAccountHeader(), accountId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraUserIdHeader(), authzInfo.UserId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraRoleHeader(), authzInfo.RoleId)

		logrus.WithFields(logrus.Fields{
			LOG_PARAM_ACCOUNT_ID: authzInfo.AccountId,
			LOG_PARAM_ROLE_ID:    authzInfo.RoleId,
			LOG_PARAM_USER_ID:    authzInfo.UserId,
		}).Debug(fmt.Sprintf("user [%s], accessing tenant [%s], with role id [%s]", authzInfo.UserId, accountId, authzInfo.RoleId))
	}

	h.Proxy.ServeHTTP(res, req)

}
