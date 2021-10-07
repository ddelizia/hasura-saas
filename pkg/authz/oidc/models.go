package oidc

import "github.com/dgrijalva/jwt-go"

type HasuraSaasClaims struct {
	HasuraSaasRole    string `json:"role,omitempty"`
	HasuraSaasAccount string `json:"account,omitempty"`
}

type JwtTokenWithClaims struct {
	*jwt.StandardClaims
	HasuraSaasNamespace HasuraSaasClaims `json:"hasura_saas_claims,omitempty"`
}
