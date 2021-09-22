package oidc

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

type UserGetterOidc struct{}

/*
GetUser Validates and JWT token and returns the decoded information of the token
*/
func (s *UserGetterOidc) GetUser(r *http.Request) (string, error) {

	jwtB64 := GetAuthHeader(r)

	if jwtB64 == "" {
		logrus.Debug("request has no authorization info, will set as anonymous user")
		return authz.ConfigAnonymousUser(), nil
	}

	parsedToken, err := ValidateJwtFunc(jwtB64, r)
	if err != nil {
		return "", errorx.IllegalState.Wrap(err, "not able to validate user as token is not valid")
	}

	return parsedToken.Subject, nil
}
