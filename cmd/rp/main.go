package main

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/rp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.IntLogger()

	sdkService := &gqlsdk.Client{
		Client: gqlsdk.NewClientBuilder(nil).BuildClient(true, gqlsdk.WithAdminRole()),
	}

	authzSvc := authz.NewService(sdkService)

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
