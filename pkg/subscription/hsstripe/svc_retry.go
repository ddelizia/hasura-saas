package hsstripe

import (
	"context"
	"fmt"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/paymentmethod"
)

//////////////////////////////////////
// Interface
//////////////////////////////////////

type StripeRetryer interface {
	Retry(ctx context.Context, input *model.RetryInput) (*model.RetryOutput, error)
}

//////////////////////////////////////
// Struct
//////////////////////////////////////

type StripeRetry struct {
	GqlreqSvc gqlreq.Service
	GqlsdkSvc gqlsdk.Service
}

//////////////////////////////////////
// Mock
//////////////////////////////////////

type StripeRetryMock struct {
	mock.Mock
}

func (m *StripeInitMock) Retry(ctx context.Context, input *model.RetryInput) (*model.RetryOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.RetryOutput), args.Error(1)
}

//////////////////////////////////////
// New
//////////////////////////////////////

func NewStripeRetry(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) StripeRetryer {
	return &StripeRetry{
		GqlreqSvc: gqlreqSvc,
		GqlsdkSvc: gqlsdkSvc,
	}
}

//////////////////////////////////////
// Method implementation
//////////////////////////////////////

// Init the subscription
func (s *StripeRetry) Retry(ctx context.Context, input *model.RetryInput) (*model.RetryOutput, error) {

	logrus.Debug("getting saas account information from hasura")
	accountInfoForCreatingSubscription, err := getStripeCustomerIdFromHasura(ctx, s.GqlsdkSvc, input.IDAccount)
	if err != nil {
		return nil, err
	}

	logrus.Debug("executing paymentId change on the provider")
	ser, err := s.retrySubscriptionOnStripe(
		ctx,
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionCustomer.StripeCustomer,
		*input.IDPaymentMethod,
	)
	if err != nil {
		return nil, err
	}

	logrus.Debug("update subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(ctx, s.GqlsdkSvc, input.IDAccount, ser)
	if err != nil {
		return nil, err
	}

	logrus.Debug("building rery response")
	return &model.RetryOutput{
		IDAccount: input.IDAccount,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}, nil

}

func (s *StripeRetry) retrySubscriptionOnStripe(ctx context.Context, c string, paymentMethodId string) (*stripe.Subscription, error) {
	// Get latest invoice for customer
	prevI := invoice.List(&stripe.InvoiceListParams{
		Customer: stripe.String(c),
	})
	var prevIn *stripe.Invoice
	if prevI.Next() {
		prevIn = prevI.Invoice()
	} else {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			model.LOG_PARAM_CUSTOMER_ID: c,
		}).Error("no invoice found")
		return nil, errorx.InternalError.New("not able to find the last invoice")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(prevIn),
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("retrieved invoice")

	// Attach PaymentMethod
	pm, err := paymentmethod.Attach(
		paymentMethodId,
		&stripe.PaymentMethodAttachParams{
			Customer: stripe.String(c),
		},
	)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error(fmt.Sprintf("paymentmethod.Attach: %v %s", err, pm.ID))
		return nil, errorx.InternalError.Wrap(err, "unable to execute retry request to the payment provider")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(pm),
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("payment attached")

	// Update invoice settings default
	updatedCustomer, err := StripeUpdateCustomerFunc(
		c,
		&stripe.CustomerParams{
			InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
				DefaultPaymentMethod: stripe.String(pm.ID),
			},
		},
	)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("unable to update customer invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to update invoice settings on the provider")
	}
	logrus.WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(updatedCustomer),
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("default payment method for customer updated")

	// Retrieve Invoice
	invoiceParams := &stripe.InvoiceParams{}
	invoiceParams.AddExpand("payment_intent")
	in, err := StripeGetInvoiceFunc(
		prevIn.ID,
		invoiceParams,
	)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("unable to retrieve invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to retrive the invoice from the payment provider")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(in),
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("invoice id updated")

	// Pay invoice Invoice
	inPay, err := StripePayInvoiceFunc(
		prevIn.ID,
		&stripe.InvoicePayParams{},
	)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("unable to pay invoice")
		return nil, errorx.InternalError.Wrap(err, "unable to invoice")
	}
	logrus.WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(inPay),
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("invoice paid")

	// Getting subscription
	ser, err := StripeGetSubFunc(in.Subscription.ID, &stripe.SubscriptionParams{})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("customer", c).Error("unable to get subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to get subscription")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(ser),
		model.LOG_PARAM_CUSTOMER_ID:     c,
		model.LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		model.LOG_PARAM_INVOICE_ID:      prevIn.ID,
	}).Debug("subscription found")

	return ser, nil
}
