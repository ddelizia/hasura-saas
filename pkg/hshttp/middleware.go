package hshttp

import (
	"net/http"
)

func MiddlewareSetContentTypeApplicationJson(next http.HandlerFunc, arguments ...interface{}) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetHeaderOnResponse(w, "Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
