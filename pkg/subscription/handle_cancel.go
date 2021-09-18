package subscription

import (
	"context"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

type CancelHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewCancelHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &CancelHandler{
		SdkSvc:     sdkSvc,
		GraphqlSvc: graphqlSvc,
	}
}

type ActionPayloadCancel struct {
	hasura.BasePayload
}

func (h *CancelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("cancel subscription request")
	actionPayload := &ActionPayloadCancel{}
	err := hshttp.GetBody(r, actionPayload)
	if err != nil {
		hshttp.WriteError(w, errorx.IllegalArgument.Wrap(err, "invalid payload for create subscription"))
		return
	}

	logrus.Debug("getting authorization information")
	authzInfo, err := h.GraphqlSvc.GetSessionInfo(actionPayload.SessionVariables)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve authz information"))
		return
	}

	logrus.Debug("retrieving subscription id on hasura")
	subscriptionId, err := h.getStripeSubscriptionID(r.Context(), authzInfo.AccountId)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "subscription could not be found"))
		return
	}

	logrus.Debug("cancelling stripe subscription")
	ser, err := h.cancelStripeSubscription(authzInfo.AccountId, subscriptionId)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "error while processing with the payment provider"))
		return
	}

	logrus.Debug("updating hasura with the cancellation status")
	_, err = updateHasuraSubscription(r.Context(), h.SdkSvc, authzInfo.AccountId, ser)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "error while updating subscription information on hasura"))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.CancelSubscriptionOutput{
		Status: string(ser.Status),
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}

/*
Retrieving the stripe subscription id from hasura
*/
func (h *CancelHandler) getStripeSubscriptionID(ctx context.Context, accountID string) (string, error) {
	result, err := h.SdkSvc.GetStripeSubscription(ctx, accountID)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_ACCOUNT_ID, accountID).Error("not able to execute GetStripeSubscription")
		return "", errorx.InternalError.Wrap(err, "not able to execute GetStripeSubscription")
	}

	if len(result.SubscriptionStatus) != 1 {
		logrus.WithField(LOG_PARAM_RESULT_LENGTH, len(result.SubscriptionStatus)).Error("hasura GetStripeSubscription has returned wrong amount of results")
		return "", errorx.InternalError.Wrap(err, "subscription not able to cancelled, contact us")
	}

	logrus.WithField(LOG_PARAM_HASURA_RESPONSE, logger.PrintStruct(result)).Debug("GetStripeSubscription result")

	return *result.SubscriptionStatus[0].StripeSubscriptionID, nil
}

/*
Requesting stripe to cancel the subscription
*/
func (h *CancelHandler) cancelStripeSubscription(accountID string, stripeSubscriptioId string) (*stripe.Subscription, error) {
	subscriptionCancel, err := sub.Cancel(stripeSubscriptioId, nil)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			LOG_PARAM_ACCOUNT_ID:      accountID,
			LOG_PARAM_SUBSCRIPTION_ID: stripeSubscriptioId,
		}).Error("unable to cancel subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to process cancel request with the provider")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(subscriptionCancel),
		LOG_PARAM_ACCOUNT_ID:      accountID,
		LOG_PARAM_SUBSCRIPTION_ID: stripeSubscriptioId,
	}).Debug("subscription done")

	return subscriptionCancel, nil
}
