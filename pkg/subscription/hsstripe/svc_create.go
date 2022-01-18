package hsstripe

import (
	"context"
	"time"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hserrorx"
	"github.com/ddelizia/hasura-saas/pkg/hstype"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/sub"
)

//////////////////////////////////////
// Interface
//////////////////////////////////////

type StripeCreator interface {
	Create(ctx context.Context, input *model.CreateInput) (*model.CreateOutput, error)
}

//////////////////////////////////////
// Struct
//////////////////////////////////////

type StripeCreate struct {
	GqlreqSvc gqlreq.Service
	GqlsdkSvc gqlsdk.Service
	// wrap stripe function call for testing purpuses
	StripeAttachPaymentmethodFunc func(id string, params *stripe.PaymentMethodAttachParams) (*stripe.PaymentMethod, error)
	StripeNewSubFunc              func(params *stripe.SubscriptionParams) (*stripe.Subscription, error)
}

//////////////////////////////////////
// Mock
//////////////////////////////////////

type StripeCreateMock struct {
	mock.Mock
}

func (m *StripeCreateMock) Create(ctx context.Context, input *model.CreateInput) (*model.CreateOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.CreateOutput), args.Error(1)
}

//////////////////////////////////////
// New
//////////////////////////////////////

func NewStripeCreate(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) StripeCreator {
	return &StripeCreate{
		GqlreqSvc:                     gqlreqSvc,
		GqlsdkSvc:                     gqlsdkSvc,
		StripeAttachPaymentmethodFunc: paymentmethod.Attach,
		StripeNewSubFunc:              sub.New,
	}
}

//////////////////////////////////////
// Method implementation
//////////////////////////////////////

// Init the subscription
func (s *StripeCreate) Create(ctx context.Context, input *model.CreateInput) (*model.CreateOutput, error) {

	logrus.WithContext(ctx).Debug("getting saas account information from hasura")
	accountInfoForCreatingSubscription, err := s.getStripeCustomerIdFromHasura(ctx, s.GqlsdkSvc, input.IDAccount)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("create subscription on stripe")
	ser, err := s.attachPaymentMethodToStripeCustomer(
		ctx,
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionCustomer.StripeCustomer,
		&input.IDPaymentMethod,
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionStatus.SubscriptionPlan.StripeCode,
		accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionStatus.SubscriptionPlan.TrialDays,
	)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("create subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(ctx, s.GqlsdkSvc, input.IDAccount, ser)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("building create response")
	return &model.CreateOutput{
		IDAccount: updatedStatus.UpdateSubscriptionStatus.Returning[0].IDAccount,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}, nil

}

func (s *StripeCreate) getStripeCustomerIdFromHasura(ctx context.Context, sdkSvc gqlsdk.Service, accountID string) (*gqlsdk.QueryGetAccountInfoForCreatingSubscription, error) {

	result, err := sdkSvc.GetAccountInfoForCreatingSubscription(ctx, accountID)
	if err != nil {
		logrus.WithError(err).Error("not able to get customer information from hasure")
		return nil, errorx.InternalError.Wrap(err, "not able to get customer information")
	}

	err = hserrorx.AssertTrue(len(result.SaasAccount) == 1,
		hserrorx.Fields{
			model.LOG_PARAM_ACCOUNT_ID: accountID,
		},
		"account not found", hstype.NewString("hasura GetAccountInfoForCreatingSubscription has returned wrong amount of results"),
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
 Create subscription to Stripe payment provider
 Steps:
 * Attach payment method to the customer
 * Update customer invoice settings with the default payment method
 * Create subscription to the plan
*/
func (s *StripeCreate) attachPaymentMethodToStripeCustomer(ctx context.Context, c string, paymentMethodId, priceId *string, trialDays *int64) (*stripe.Subscription, error) {

	if priceId == nil {
		logrus.WithContext(ctx).WithField(model.LOG_PARAM_STRIPE_PLAN_ID, priceId).Error("stripe price does not exists")
		return nil, errorx.InternalError.New("stripe price does not exists")
	}

	if paymentMethodId != nil {
		// Attach payment method to the customer
		params := &stripe.PaymentMethodAttachParams{
			Customer: hstype.NewString(c),
		}
		pm, err := s.StripeAttachPaymentmethodFunc(
			*paymentMethodId,
			params,
		)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("payment attachment failed")
			return nil, errorx.InternalError.Wrap(err, "unable to attach payment to the subscription")
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(pm),
			model.LOG_PARAM_CUSTOMER_ID:     c,
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
			logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("unable to update customer invoice")
			return nil, errorx.InternalError.Wrap(err, "unable to update invoice settings")
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(updatedCustomer),
			model.LOG_PARAM_CUSTOMER_ID:     c,
		}).Debug("default payment method for customer updated")

	}

	// Create subscription to the plan
	var trialEnd *int64 = nil
	if trialDays != nil {
		trialEnd = stripe.Int64(time.Now().AddDate(0, 0, int(*trialDays)).Unix())
		logrus.WithContext(ctx).Debug("trial ends: ", trialEnd)
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(c),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: stripe.String(*priceId),
			},
		},
		TrialEnd: trialEnd,
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	ser, err := s.StripeNewSubFunc(subscriptionParams)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_CUSTOMER_ID, c).Error("unable to create subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to create subscription")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(ser),
		model.LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		model.LOG_PARAM_CUSTOMER_ID:     c,
	}).Debug("subscription done")

	return ser, nil
}

