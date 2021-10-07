package oidc

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

func getJwtB64PublicKey() (*keyfunc.JWKs, error) {
	// Create the keyfunc options. Refresh the JWKS every hour and log errors.
	refreshInterval := time.Hour
	options := keyfunc.Options{
		RefreshInterval: &refreshInterval,
		RefreshErrorHandler: func(err error) {
			logrus.WithError(err).Error("there was an error with the jwt.KeyFunc")
		},
	}

	// Create the JWKS from the resource at the given URL.
	return keyfunc.Get(ConfigJwksUrl(), options)
}

var ValidateJwtFunc = validateJwt

func validateJwt(jwtB64 string, r *http.Request) (*JwtTokenWithClaims, error) {

	parsedToken := &JwtTokenWithClaims{}

	// trying to get token from header cache value
	decodedHeader := r.Header.Get(ConfigHeaderNameDecodedJwt())
	if decodedHeader != "" {
		if err := json.Unmarshal([]byte(decodedHeader), parsedToken); err != nil {
			logrus.WithError(err).Error("not able to parse existing decoded header")
			return nil, errorx.InternalError.Wrap(err, "not able to parse existing decoded header")
		}

		return parsedToken, nil
	}

	jwks, err := getJwtB64PublicKey()
	if err != nil {
		return nil, errorx.ExternalError.New("not able to retrieve jwks")
	}

	token, err := jwt.ParseWithClaims(jwtB64, parsedToken, jwks.KeyFunc)
	if err != nil || !token.Valid {
		logrus.WithError(err).WithField("jwt", jwtB64).Error("not able to parse token")
		return nil, errorx.IllegalArgument.New("token is invalid")
	}

	// storing token into header
	data, err := json.Marshal(parsedToken)
	if err != nil {
		logrus.WithError(err).Error("not able to marshal token")
		return nil, errorx.InternalError.Wrap(err, "not able to marshal token")
	}

	hshttp.SetHaderOnRequest(r, ConfigHeaderNameDecodedJwt(), string(data))

	return parsedToken, nil
}
