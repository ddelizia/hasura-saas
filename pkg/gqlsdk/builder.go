package gqlsdk

import (
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
)

type ClientBuilder interface {
	BuildClient(setHasuraAdminSecret bool, additional ...clientv2.RequestInterceptor) *clientv2.Client
}

type ClientBuilderImpl struct {
	HttpClient *http.Client
}

func NewClientBuilder(httpClient *http.Client) ClientBuilder {
	return &ClientBuilderImpl{
		HttpClient: httpClient,
	}
}

func (b *ClientBuilderImpl) BuildClient(setHasuraAdminSecret bool, additional ...clientv2.RequestInterceptor) *clientv2.Client {
	if b.HttpClient == nil {
		b.HttpClient = http.DefaultClient
	}

	var interceptors []clientv2.RequestInterceptor
	if setHasuraAdminSecret {
		interceptors = append(interceptors, HasuraAdminInterceptor())
	}
	interceptors = append(interceptors, additional...)

	return clientv2.NewClient(
		b.HttpClient,
		gqlreq.ConfigGraphqlURL(),
		clientv2.ChainInterceptor(interceptors...),
	)
}
