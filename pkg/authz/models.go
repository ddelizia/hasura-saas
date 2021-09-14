package authz

import "github.com/dgrijalva/jwt-go"

type SaasMembershipResponse struct {
	SaasMembership []SaasMembership `json:"saas_membership"`
}

type SaasMembership struct {
	SaasRole SaasRole `json:"saas_role"`
}

type SaasRole struct {
	HasuraRole string `json:"hasura_role"`
}

type JwtTokenWithClaims struct {
	*jwt.StandardClaims
	TokenType string
}

type AuthzInfo struct {
	RoleId    string
	UserId    string
	AccountId string
	Jwt       string
}
