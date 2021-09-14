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

	logrus.WithField("event", event.Type).Debug("storing event")
	var data map[string]interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		logrus.WithError(err).WithField("body", string(b)).Error("not able to unmarshal")
	}
	h.SdkSvc.AddSubscriptionEvent(r.Context(), event.Type, data)

	switch event.Type {
	case "payment_intent.succeeded":

	case "checkout.session.completed":
		// Payment is successful and the subscription is created.
		// You should provision the subscription.
	case "invoice.paid":
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
	case "invoice.payment_failed":
		// The payment failed or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
	case "customer.created":
		// Customer has been created
	default:
		// unhandled event type
		logrus.WithField("eventType", event.Type).Info("event not processed")
	}

	w.WriteHeader(http.StatusOK)
}
