package gqlreq

import (
	"net/http"

	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type HeaderInfoGetter interface {
	GetHeaderInfo(r *http.Request) (*HeaderInfo, error)
}

type HeaderInfoGetterImpl struct {
}

type HeaderInfoGetterMock struct {
	mock.Mock
}

func (m *HeaderInfoGetterMock) GetHeaderInfo(r *http.Request) (*HeaderInfo, error) {
	args := m.Called(r)
	return args.Get(0).(*HeaderInfo), args.Error(1)
}

func (s *HeaderInfoGetterImpl) GetHeaderInfo(r *http.Request) (*HeaderInfo, error) {
	accountId := r.Header.Get(ConfigHasuraAccountHeader())
	userId := r.Header.Get(ConfigHasuraAccountHeader())
	role := r.Header.Get(ConfigHasuraAccountHeader())

	logrus.WithFields(logrus.Fields{
		"accountId": accountId,
		"userId":    userId,
		"role":      role,
	}).Debug("authz headers set")

	if accountId == "" || userId == "" || role == "" {
		return nil, errorx.IllegalArgument.New("missing parametes: [accountId] [userId] [role]")
	}

	return &HeaderInfo{
		AccountId: accountId,
		UserId:    userId,
		Role:      role,
	}, nil
}
