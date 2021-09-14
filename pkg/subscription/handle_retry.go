package subscription

import (
	"fmt"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/sub"
)

type RetryHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewRetryHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &RetryHandler{
		GraphqlSvc: graphqlSvc,
		SdkSvc:     sdkSvc,
	}
}

type ActionPayloadRetry struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.RetrySubscriptionInput `json:"data"`
	} `json:"input"`
}

func (h *RetryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("retry subscription request")
	actionPayload := &ActionPayloadRetry{}
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

	logrus.Debug("executing paymentId change on the provider")
	ser, err := h.retrySubscriptionOnStripe(
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionCustomer.StripeCustomer,
		actionPayload.Input.Data.PaymentMethodID,
		"", //TODO invoice should be placed here
	)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to assign the new payment method with the payment provider"))
		return
	}

	logrus.Debug("update subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(r.Context(), h.SdkSvc, authzInfo.AccountId, ser)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to store subscription information"))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.RetrySubscriptionOutput{
		AccountID: authzInfo.AccountId,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}

func (h *RetryHandler) retrySubscriptionOnStripe(c string, paymentMethodId string, invoiceId string) (*stripe.Subscription, error) {
	// Attach PaymentMethod
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(c),
	}
	pm, err := paymentmethod.Attach(
		paymentMethodId,
		params,
	)
	if err != nil {
		logrus.WithError(err).WithField("customer", c).Error(fmt.Sprintf("paymentmethod.Attach: %v %s", err, pm.ID))
		return nil, errorx.InternalError.Wrap(err, "unable to execute retry request to the payment provider")
	}
	logrus.WithFields(logrus.Fields{
		"stripe.attach": logger.PrintStruct(pm),
		"customer":      c,
	}).Debug("payment attached")

	// Update invoice settings default
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
		logrus.WithError(err).WithField("customer", c).Error("unable to update customer invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to update invoice settings on the provider")
	}
	logrus.WithFields(logrus.Fields{
		"stripe.customer.updated": logger.PrintStruct(updatedCustomer),
		"customer":                c,
	}).Debug("default payment method for customer updated")

	// Retrieve Invoice
	invoiceParams := &stripe.InvoiceParams{}
	invoiceParams.AddExpand("payment_intent")
	in, err := invoice.Get(
		invoiceId,
		invoiceParams,
	)
	if err != nil {
		logrus.WithError(err).WithField("customer", c).Error("unable to retrieve invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to retrive the invoice from the payment provider")
	}
	logrus.WithFields(logrus.Fields{
		"stripe.invoice": logger.PrintStruct(in),
		"customer":       c,
	}).Debug("invoice id updated")

	// Getting subscription
	subscriptionParams := &stripe.SubscriptionParams{}
	ser, err := sub.Get(in.Subscription.ID, subscriptionParams)
	if err != nil {
		logrus.WithError(err).WithField("customer", c).Error("unable to get subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to get subscription")
	}
	logrus.WithFields(logrus.Fields{
		"stripe.subscription": logger.PrintStruct(ser),
		"customer":            c,
		"subscriptionId":      ser.ID,
	}).Debug("subscription found")

	return ser, nil
}
