package oidc

import (
	"net/http"
	"strings"
)

func GetAuthHeader(r *http.Request) string {

	authH := r.Header.Get("Authorization")

	return strings.TrimPrefix(authH, "Bearer ")
}
