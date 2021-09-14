package gqlreq

import (
	"strings"

	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type SessionInfoGetter interface {
	GetSessionInfo(r map[string]interface{}) (*HeaderInfo, error)
}

type SessionInfoGetterImpl struct {
}

type SessionInfoGetterMock struct {
	mock.Mock
}

func (m *SessionInfoGetterMock) GetSessionInfo(r map[string]interface{}) (*HeaderInfo, error) {
	args := m.Called(r)
	return args.Get(0).(*HeaderInfo), args.Error(1)
}

func (s *HeaderInfoGetterImpl) GetSessionInfo(r map[string]interface{}) (*HeaderInfo, error) {

	// TODO ignore case is needed here
	accountId := r[strings.ToLower(ConfigHasuraAccountHeader())]
	userId := r[strings.ToLower(ConfigHasuraUserIdHeader())]
	role := r[strings.ToLower(ConfigHasuraRoleHeader())]

	logrus.WithFields(logrus.Fields{
		"accountId": accountId,
		"userId":    userId,
		"role":      role,
	}).Debug("authz session set")

	if accountId == "" || userId == "" || role == "" {
		return nil, errorx.IllegalArgument.New("missing parametes: [accountId] [userId] [role]")
	}

	return &HeaderInfo{
		AccountId: accountId.(string),
		UserId:    userId.(string),
		Role:      role.(string),
	}, nil
}
