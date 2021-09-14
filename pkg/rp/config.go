package rp

import (
	"net/url"

	"github.com/ddelizia/hasura-saas/pkg/env"
)

func ConfigListenAddress() string {
	return env.GetString("rp.server.listenAddress")
}

/*
ES configuration
*/

func ConfigEsUrl() *url.URL {
	return env.GetUrl("rp.es.url")
}

func ConfigEsPublicIndex() string {
	return env.GetString("rp.es.index.public")
}

func ConfigEsPrivateIndex() string {
	return env.GetString("rp.es.index.private")
}

func ConfigEsAuthorizationHeader() string {
	return env.GetString("rp.es.headerNames.authorization")
}

/*
HASURA configuration
*/

func ConfigHasuraUrl() *url.URL {
	return env.GetUrl("rp.graphql.url")
}
