package oidc

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

type RoleGetterOidc struct{}

/*
GetRole retrieves the user role from jwt
*/
func (s *RoleGetterOidc) GetRole(r *http.Request) (string, error) {

	jwtB64 := GetAuthHeader(r)

	if jwtB64 == "" {
		logrus.Debug("request has no authorization info, will set as anonymous role")
		return authz.ConfigAnonymousRole(), nil
	}

	parsedToken, err := ValidateJwtFunc(jwtB64, r)
	if err != nil {
		return "", errorx.IllegalState.Wrap(err, "not able to validate user as token is not valid")
	}

	return parsedToken.HasuraSaasNamespace.HasuraSaasRole, nil
}
