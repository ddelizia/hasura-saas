package main

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hshttp/hsmiddleware"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription"
	"github.com/ddelizia/hasura-saas/pkg/subscription/hsstripe"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

func main() {
	logger.IntLogger()
	logrus.WithField("serverPort", subscription.ConfigWebhookListenAddress()).Info("starting subscription server")

	stripe.Key = subscription.ConfigApiKey()
	stripe.DefaultLeveledLogger = logrus.New().WithField("component", "stripe")

	graphqlSevice := gqlreq.NewService()
	sdkService := &gqlsdk.Client{
		Client: gqlsdk.NewClientBuilder(nil).BuildClient(true, gqlsdk.WithAdminRole()),
	}
	hsStripeService := hsstripe.NewService(graphqlSevice, sdkService)

	r := mux.NewRouter()

	handlerInit := subscription.NewInitHandler(hsStripeService)
	r.Handle("/init",
		hsmiddleware.Chain(
			handlerInit.ServeHTTP,
			hsmiddleware.LogRequest(),
			hsmiddleware.ActionBodyToContext(&subscription.ActionPayloadInit{}),
			hsmiddleware.AuthzFromSession(graphqlSevice),
			hsmiddleware.Json(),
		)).Methods("POST")

	handlerCreate := subscription.NewCreateHandler(graphqlSevice, sdkService)
	r.Handle("/create",
		hsmiddleware.Chain(
			handlerCreate.ServeHTTP,
			hsmiddleware.LogRequest(),
			hsmiddleware.Json(),
		)).Methods("POST")

	handlerCancel := subscription.NewCancelHandler(graphqlSevice, sdkService)
	r.Handle("/cancel",
		hsmiddleware.Chain(
			handlerCancel.ServeHTTP,
			hsmiddleware.LogRequest(),
			hsmiddleware.Json(),
		)).Methods("POST")

	handlerRetry := subscription.NewRetryHandler(graphqlSevice, sdkService)
	r.Handle("/retry",
		hsmiddleware.Chain(
			handlerRetry.ServeHTTP,
			hsmiddleware.LogRequest(),
			hsmiddleware.Json(),
		)).Methods("POST")

	handlerChange := subscription.NewChangeHandler(graphqlSevice, sdkService)
	r.Handle("/change",
		hsmiddleware.Chain(
			handlerChange.ServeHTTP,
			hsmiddleware.LogRequest(),
			hsmiddleware.Json(),
		)).Methods("POST")

	handlerWebhook := subscription.NewWebhookHandler(sdkService)
	r.Handle("/webhook", handlerWebhook).Methods("POST")

	http.Handle("/", r)

	if err := http.ListenAndServe(subscription.ConfigWebhookListenAddress(), nil); err != nil {
		panic(err)
	}
}
