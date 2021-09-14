package authz

import (
	"context"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/mock"
)

type AuthInfoGetter interface {
	GetAuthInfo(c context.Context, jwt string, account string) (*AuthzInfo, error)
	GetAuthInfoFromRequest(r *http.Request) (*AuthzInfo, error)
}

type AuthInfoGetterMock struct {
	mock.Mock
}

func (m *AuthInfoGetterMock) GetAuthInfo(c context.Context, jwtB64 string, account string) (*AuthzInfo, error) {
	args := m.Called(c, jwtB64, account)
	return args.Get(0).(*AuthzInfo), args.Error(1)
}

func (m *AuthInfoGetterMock) GetAuthInfoFromRequest(r *http.Request) (*AuthzInfo, error) {
	args := m.Called(r)
	return args.Get(0).(*AuthzInfo), args.Error(1)
}

type AuthInfoGetterImpl struct {
	UserGetter
	RoleGetter
}

// GetAuthInfo returns all the authentication information for a JWT token.
// It takes care of validating both token and membership of the user in the account (if provided)
func (s *AuthInfoGetterImpl) GetAuthInfo(c context.Context, jwtB64 string, account string) (*AuthzInfo, error) {

	user, err := s.GetUser(c, jwtB64)

	if err != nil {
		return nil, errorx.Decorate(err, "token is invalid")
	}

	role, err := s.GetRole(c, account, user)

	if err != nil {
		return nil, errorx.Decorate(err, "role cannot be calculated")
	}

	return &AuthzInfo{
		UserId:    user,
		RoleId:    role,
		AccountId: account,
		Jwt:       jwtB64,
	}, nil
}

func (s *AuthInfoGetterImpl) GetAuthInfoFromRequest(r *http.Request) (*AuthzInfo, error) {

	authToken := r.Header.Get(hshttp.ConfigJWTHeader())
	accountId := r.Header.Get(hshttp.ConfigAccountIdHeader())

	return s.GetAuthInfo(r.Context(), authToken, accountId)
}
