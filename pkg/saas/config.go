package saas

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigListenAddress() string {
	return env.GetString("saas.server.listenAddress")
}
