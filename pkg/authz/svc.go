package authz

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type Service interface {
	RoleGetter
	UserGetter
	AccountGetter
	AuthInfoGetter
}

type RoleGetter interface {
	GetRole(r *http.Request) (string, error)
}

type RoleGetterMock struct {
	RoleGetter
	mock.Mock
}

func (m *RoleGetterMock) GetRole(r *http.Request) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}

type UserGetter interface {
	GetUser(r *http.Request) (string, error)
}

type UserGetterMock struct {
	UserGetter
	mock.Mock
}

func (m *UserGetterMock) GetUser(r *http.Request) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}

type AccountGetter interface {
	GetAccount(r *http.Request) (string, error)
}

type AccountGetterMock struct {
	AccountGetter
	mock.Mock
}

func (m *AccountGetterMock) GetAccount(r *http.Request) (string, error) {
	args := m.Called(r)
	return args.String(0), args.Error(1)
}
