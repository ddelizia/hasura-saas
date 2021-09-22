package oidc

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigJwksUrl() string {
	return env.GetString("authz.oidc.jwks.url")
}

func ConfigHeaderNameDecodedJwt() string {
	return env.GetString("authz.oidc.headerNames.decodedJwt")
}