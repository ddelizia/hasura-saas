package hsmiddleware

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
)

func AuthzFromSession(gqlSvc gqlreq.Service) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionVars := hscontext.ActionSessionVariablesValue(r.Context())
			authzInfo, err := gqlSvc.GetSessionInfo(sessionVars)
			if err != nil {
				hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve authz information"))
				return
			}
			newContext := hscontext.WithAuthzInfoValue(r.Context(), authzInfo)
			r = r.WithContext(newContext)
			next.ServeHTTP(w, r)
		})
	}
}
