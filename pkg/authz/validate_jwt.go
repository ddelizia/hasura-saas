package authz

import (
	"time"

	"github.com/MicahParks/keyfunc"
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

func ValidateJwt(jwtB64 string) (*JwtTokenWithClaims, error) {

	jwks, err := getJwtB64PublicKey()
	if err != nil {
		return nil, errorx.ExternalError.New("not able to retrieve jwks")
	}

	parsedToken := &JwtTokenWithClaims{}
	token, err := jwt.ParseWithClaims(jwtB64, parsedToken, jwks.KeyFunc)
	if err != nil || !token.Valid {
		logrus.WithError(err).WithField("jwt", jwtB64).Error("not able to parse token")
		return nil, errorx.IllegalArgument.New("token is invalid")
	}

	return parsedToken, nil
}
