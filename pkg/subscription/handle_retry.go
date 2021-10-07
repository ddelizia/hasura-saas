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
		IDAccount: authzInfo.AccountId,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}

	logrus.WithFields(logrus.Fields{
		LOG_PARAM_ACCOUNT_ID:      authzInfo.AccountId,
		LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		LOG_PARAM_CUSTOMER_ID:     ser.Customer.ID,
		LOG_PARAM_PLAN_ID:         ser.Plan.ID,
	}).Info("subscription retry done")
}

func (h *RetryHandler) retrySubscriptionOnStripe(c string, paymentMethodId string) (*stripe.Subscription, error) {
	// Get latest invoice for customer
	prevI := invoice.List(&stripe.InvoiceListParams{
		Customer: stripe.String(c),
	})
	var prevIn *stripe.Invoice
	if prevI.Next() {
		prevIn = prevI.Invoice()
	} else {
		logrus.WithFields(logrus.Fields{
			LOG_PARAM_CUSTOMER_ID: c,
		}).Error("no invoice found")
		return nil, errorx.InternalError.New("not able to find the last invoice")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(prevIn),
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("retrieved invoice")

	// Attach PaymentMethod
	pm, err := paymentmethod.Attach(
		paymentMethodId,
		&stripe.PaymentMethodAttachParams{
			Customer: stripe.String(c),
		},
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error(fmt.Sprintf("paymentmethod.Attach: %v %s", err, pm.ID))
		return nil, errorx.InternalError.Wrap(err, "unable to execute retry request to the payment provider")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(pm),
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("payment attached")

	// Update invoice settings default
	updatedCustomer, err := customer.Update(
		c,
		&stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(pm.ID),
			},
		},
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("unable to update customer invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to update invoice settings on the provider")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(updatedCustomer),
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("default payment method for customer updated")

	// Retrieve Invoice
	invoiceParams := &stripe.InvoiceParams{}
	invoiceParams.AddExpand("payment_intent")
	in, err := invoice.Get(
		prevIn.ID,
		invoiceParams,
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("unable to retrieve invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to retrive the invoice from the payment provider")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(in),
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("invoice id updated")

	// Pay invoice Invoice
	inPay, err := invoice.Pay(
		prevIn.ID,
		&stripe.InvoicePayParams{},
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c).Error("unable to pay invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to invoice")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(inPay),
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
		LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("invoice paid")

	// Getting subscription
	ser, err := sub.Get(in.Subscription.ID, &stripe.SubscriptionParams{})
	if err != nil {
		logrus.WithError(err).WithField("customer", c).Error("unable to get subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to get subscription")
	}
	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(ser),
		LOG_PARAM_CUSTOMER_ID:     c,
		LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		LOG_PARAM_INVOICE_ID:      prevIn.ID,
	}).Debug("subscription found")

	return ser, nil
}
