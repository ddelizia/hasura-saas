package subscription

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/hstype"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/sub"
)

type CreateHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewCreateHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &CreateHandler{
		SdkSvc:     sdkSvc,
		GraphqlSvc: graphqlSvc,
	}
}

type ActionCreatePayload struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.CreateSubscriptionInput `json:"data"`
	} `json:"input"`
}

/*
Handle subscription creation
*/
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("create subscription request")
	actionPayload := &ActionCreatePayload{}
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

	logrus.Debug("getting saas account information from hasura")
	accountInfoForCreatingSubscription, err := getStripeCustomerIdFromHasura(r.Context(), h.SdkSvc, authzInfo)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve customer"))
		return
	}

	logrus.Debug("create subscription on stripe")
	ser, err := attachPaymentMethodToStripeCustomer(
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionCustomer.StripeCustomer,
		actionPayload.Input.Data.PaymentMethodID,
		*accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionStatus.SubscriptionPlan.StripeCode,
	)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable create subscription on payment provider"))
		return
	}

	logrus.Debug("create subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(r.Context(), h.SdkSvc, authzInfo.AccountId, ser)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to store subscription information"))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.CreateSubscriptionOutput{
		AccountID: updatedStatus.UpdateSubscriptionStatus.Returning[0].IDAccount,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}

/*
 Create subscription to Stripe payment provider
 Steps:
 * Attach payment method to the customer
 * Update customer invoice settings with the default payment method
 * Create subscription to the plan
*/
func attachPaymentMethodToStripeCustomer(c string, paymentMethodId string, priceId string) (*stripe.Subscription, error) {
	// Attach payment method to the customer
	params := &stripe.PaymentMethodAttachParams{
		Customer: hstype.NewString(c),
	}
	pm, err := paymentmethod.Attach(
		paymentMethodId,
		params,
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("payment attachment failde")
		return nil, errorx.InternalError.Wrap(err, "unable to attach payment to the subscription")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(pm),
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("payment attached")

	// Update customer invoice settings with the default payment method
	customerParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm.ID),
		},
	}
	updatedCustomer, err := customer.Update(
		c,
		customerParams,
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("unable to update customer invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to update invoice settings")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(updatedCustomer),
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("default payment method for customer updated")

	// Create subscription to the plan
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(c),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(priceId),
			},
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	ser, err := sub.New(subscriptionParams)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("unable to create subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to create subscription")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(ser),
		LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("subscription done")

	return ser, nil
}
