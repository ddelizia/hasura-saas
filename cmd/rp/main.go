package main

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz/oidc"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/rp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.IntLogger()

	authzSvc := oidc.NewService()

	proxyHasura := rp.NewHasuraService(authzSvc)

	logrus.WithField("serverPort", rp.ConfigListenAddress()).Info("starting reverse proxy server")

	rtr := mux.NewRouter()

	// Graphql proxy
	rtr.Handle("/graphql", proxyHasura)

	http.Handle("/", rtr)

	if err := http.ListenAndServe(rp.ConfigListenAddress(), nil); err != nil {
		panic(err)
	}
}
