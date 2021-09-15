package authz

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type RoleGetter interface {
	GetRole(c context.Context, account string, user string) (string, error)
}

type RoleGetterMock struct {
	mock.Mock
}

func (m *RoleGetterMock) GetRole(c context.Context, account string, user string) (string, error) {
	args := m.Called(c, account, user)
	return args.String(0), args.Error(1)
}

type RoleGetterImpl struct {
	GraphQlSvc gqlsdk.Service
}

/*
GetRole retrieves the user role in the membership Hasura table
*/
func (s *RoleGetterImpl) GetRole(c context.Context, account string, user string) (string, error) {
	if user == "" {
		return ConfigAnonymousRole(), nil
	}

	if account == "" {
		return ConfigLoggedInRole(), nil
	}

	response, err := s.GraphQlSvc.GetRoleForUserAndAccount(c, user, account)
	
	if err != nil {
		return "", errorx.ExternalError.New("error executing the request to graphql hasura")
	}
	if len(response.SaasMembership) == 0 {
		logrus.WithField("user", user).Debug("use not found, setting default logged in role")
		return ConfigLoggedInRole(), nil
	}
	if len(response.SaasMembership) > 1 {
		logrus.WithField("user", user).Debug("more than one user has been found with the same name in hasura")
		return "", errorx.IllegalState.New("more than one user has been found")

	}

	return response.SaasMembership[0].IDRole, nil
}
