package hsmiddleware

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/hshttp"
)

func Json() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hshttp.SetHeaderOnResponse(w, "Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
