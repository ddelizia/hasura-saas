package authz

import (
	"context"

	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

var ValidateJwtFunc = ValidateJwt

type UserGetter interface {
	GetUser(c context.Context, jwtB64 string) (string, error)
}

type UserGetterMock struct {
	mock.Mock
}

func (m *UserGetterMock) GetUser(c context.Context, jwtB64 string) (string, error) {
	args := m.Called(c, jwtB64)
	return args.String(0), args.Error(1)
}

type UserGetterImpl struct{}

/*
GetUser Validates and JWT token and returns the decoded information of the token
*/
func (s *UserGetterImpl) GetUser(c context.Context, jwtB64 string) (string, error) {

	if jwtB64 == "" {
		logrus.Debug("request has no authorization info, will set as anonymous user")
		return ConfigAnonymousUser(), nil
	}

	parsedToken, err := ValidateJwtFunc(jwtB64)
	if err != nil {
		return "", errorx.IllegalState.Wrap(err, "not able to validate user as token is not valid")
	}

	return parsedToken.Subject, nil
}
