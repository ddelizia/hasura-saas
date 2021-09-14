package gqlreq

import (
	"net/http"

	"github.com/machinebox/graphql"
)

type Service interface {
	Executer
	HeaderInfoGetter
	SessionInfoGetter
}

type service struct {
	ExecuterImpl
	HeaderInfoGetterImpl
	SessionInfoGetterImpl
}

func NewService() Service {
	client := graphql.NewClient(ConfigGraphqlURL(), graphql.WithHTTPClient(http.DefaultClient))
	return &service{
		ExecuterImpl: ExecuterImpl{Client: client},
	}
}
