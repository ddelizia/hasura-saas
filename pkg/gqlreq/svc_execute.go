package gqlreq

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type Executer interface {
	Execute(c context.Context, q string, headers []RequestHeader, vars []RequestVar, setHasuraAdminSecret bool, result interface{}) error
}

type ExecuterMock struct {
	mock.Mock
}

func (m *ExecuterMock) Execute(c context.Context, q string, headers []RequestHeader, vars []RequestVar, setHasuraAdminSecret bool, result interface{}) error {
	args := m.Called(c, q, headers, vars, setHasuraAdminSecret, result)
	return args.Error(0)
}

type GraphQlClient interface {
	Run(c context.Context, req *graphql.Request, res interface{}) error
}

type ExecuterImpl struct {
	Client GraphQlClient
}

func (s *ExecuterImpl) Execute(c context.Context, q string, headers []RequestHeader, vars []RequestVar, setHasuraAdminSecret bool, result interface{}) error {
	req := graphql.NewRequest(q)

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value)
	}

	for _, v := range vars {
		req.Var(v.Key, v.Value)
	}

	if setHasuraAdminSecret {
		req.Header.Set("X-Hasura-Admin-Secret", ConfigAdminSecret())
	}

	err := s.Client.Run(c, req, result)

	logrus.WithFields(logrus.Fields{
		"body":     q,
		"headers":  logger.PrintStruct(headers),
		"vars":     logger.PrintStruct(vars),
		"secret":   ConfigAdminSecret(),
		"response": logger.PrintStruct(result),
		"error":    logger.PrintStruct(err),
	}).Debug("executing graphql request")

	return err
}
