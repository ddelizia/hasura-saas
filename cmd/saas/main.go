package main

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hshttp/hsmiddleware"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/saas"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.IntLogger()
	logrus.WithField("serverPort", saas.ConfigListenAddress()).Info("starting saas server")

	graphqlSevice := gqlreq.NewService()
	sdkService := &gqlsdk.Client{
		Client: gqlsdk.NewClientBuilder(nil).BuildClient(true, gqlsdk.WithAdminRole()),
	}

	r := mux.NewRouter()

	handlerSetCurrentAccount := saas.NewSetCurrentAccountHandler(graphqlSevice, sdkService)
	r.Handle("/setCurrentAccount", hsmiddleware.Chain(
		handlerSetCurrentAccount.ServeHTTP,
		hsmiddleware.LogRequest(),
		hsmiddleware.Json(),
	)).Methods("POST")

	handlerGetCurrentAccount := saas.NewGetCurrentAccountHandler(graphqlSevice, sdkService)
	r.Handle("/getCurrentAccount", hsmiddleware.Chain(
		handlerGetCurrentAccount.ServeHTTP,
		hsmiddleware.LogRequest(),
		hsmiddleware.Json(),
	)).Methods("POST")

	http.Handle("/", r)

	if err := http.ListenAndServe(saas.ConfigListenAddress(), nil); err != nil {
		panic(err)
	}
}
