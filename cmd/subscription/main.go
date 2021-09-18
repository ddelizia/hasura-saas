package main

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription"
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

	r := mux.NewRouter()

	handlerInit := subscription.NewInitHandler(graphqlSevice, sdkService)
	r.Handle("/init", hshttp.MiddlewareChain(
		handlerInit.ServeHTTP,
		hshttp.MiddlewareLogRequest,
		hshttp.MiddlewareSetContentTypeApplicationJson,
	)).Methods("POST")

	handlerCreate := subscription.NewCreateHandler(graphqlSevice, sdkService)
	r.Handle("/create",
		hshttp.MiddlewareChain(
			handlerCreate.ServeHTTP,
			hshttp.MiddlewareLogRequest,
			hshttp.MiddlewareSetContentTypeApplicationJson,
		)).Methods("POST")

	handlerCancel := subscription.NewCancelHandler(graphqlSevice, sdkService)
	r.Handle("/cancel",
		hshttp.MiddlewareChain(
			handlerCancel.ServeHTTP,
			hshttp.MiddlewareLogRequest,
			hshttp.MiddlewareSetContentTypeApplicationJson,
		)).Methods("POST")

	handlerRetry := subscription.NewRetryHandler(graphqlSevice, sdkService)
	r.Handle("/retry",
		hshttp.MiddlewareChain(
			handlerRetry.ServeHTTP,
			hshttp.MiddlewareLogRequest,
			hshttp.MiddlewareSetContentTypeApplicationJson,
		)).Methods("POST")

	handlerChange := subscription.NewChangeHandler(graphqlSevice, sdkService)
	r.Handle("/change",
		hshttp.MiddlewareChain(
			handlerChange.ServeHTTP,
			hshttp.MiddlewareLogRequest,
			hshttp.MiddlewareSetContentTypeApplicationJson,
		)).Methods("POST")

	handlerWebhook := subscription.NewWebhookHandler(sdkService)
	r.Handle("/webhook", handlerWebhook).Methods("POST")

	http.Handle("/", r)

	if err := http.ListenAndServe(subscription.ConfigWebhookListenAddress(), nil); err != nil {
		panic(err)
	}
}
