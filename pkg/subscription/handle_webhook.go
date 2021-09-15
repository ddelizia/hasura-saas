package subscription

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/webhook"
)

type WebhookHandler struct {
	SdkSvc gqlsdk.Service
}

func NewWebhookHandler(sdkSvc gqlsdk.Service) http.Handler {
	return &WebhookHandler{
		SdkSvc: sdkSvc,
	}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error("not able to read webhook body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	event, err := webhook.ConstructEvent(b, r.Header.Get("Stripe-Signature"), ConfigWebhookSecret())
	if err != nil {
		logrus.
			WithError(err).
			WithFields(logrus.Fields{
				"header":    r.Header.Get("Stripe-Signature"),
				"secret":    ConfigWebhookSecret(),
				"eventbody": string(b),
			}).
			Error("not able to construct event so no type can be derived")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logrus.WithField("event", event.Type).Debug("storing event into hasura")
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		logrus.WithError(err).WithField("body", string(b)).Error("not able to unmarshal")
	}
	result, _ := h.SdkSvc.AddSubscriptionEvent(r.Context(), event.Type, data)

	EventMapping(r.Context(), event, result.InsertSubscriptionEvent.Returning[0].ID, h.SdkSvc)

	w.WriteHeader(http.StatusOK)
}
