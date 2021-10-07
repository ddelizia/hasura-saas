package authz

import (
	"net/http"

	"github.com/joomcode/errorx"
	"github.com/stretchr/testify/mock"
)

type AuthInfoGetter interface {
	GetAuthInfo(r *http.Request) (*AuthzInfo, error)
}

type AuthInfoGetterMock struct {
	mock.Mock
}

func (m *AuthInfoGetterMock) GetAuthInfo(r *http.Request) (*AuthzInfo, error) {
	args := m.Called(r)
	return args.Get(0).(*AuthzInfo), args.Error(1)
}

func (m *AuthInfoGetterMock) GetAuthInfoFromRequest(r *http.Request) (*AuthzInfo, error) {
	args := m.Called(r)
	return args.Get(0).(*AuthzInfo), args.Error(1)
}

type AuthInfoGetterImpl struct {
	UserGetter
	RoleGetter
	AccountGetter
}

// GetAuthInfo returns all the authentication information for a JWT token.
// It takes care of validating both token and membership of the user in the account (if provided)
func (s *AuthInfoGetterImpl) GetAuthInfo(r *http.Request) (*AuthzInfo, error) {

	user, err := s.GetUser(r)
	if err != nil {
		return nil, errorx.Decorate(err, "token is invalid")
	}

	role, err := s.GetRole(r)
	if err != nil {
		return nil, errorx.Decorate(err, "role cannot be calculated")
	}

	account, err := s.GetAccount(r)
	if err != nil {
		return nil, errorx.Decorate(err, "account cannot be calculated")
	}

	return &AuthzInfo{
		UserId:    user,
		RoleId:    role,
		AccountId: account,
	}, nil
}
