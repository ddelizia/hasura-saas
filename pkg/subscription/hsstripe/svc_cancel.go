package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go"
)

//////////////////////////////////////
// Interface
//////////////////////////////////////

type StripeCanceler interface {
	Cancel(ctx context.Context, input *model.CancelInput) (*model.CancelOutput, error)
}

//////////////////////////////////////
// Struct
//////////////////////////////////////

type StripeCancel struct {
	GqlreqSvc gqlreq.Service
	GqlsdkSvc gqlsdk.Service
}

//////////////////////////////////////
// Mock
//////////////////////////////////////

type StripeCancelMock struct {
	mock.Mock
}

func (m *StripeCancelMock) Cancel(ctx context.Context, input *model.CancelInput) (*model.CancelOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.CancelOutput), args.Error(1)
}

//////////////////////////////////////
// New
//////////////////////////////////////

func NewStripeCancel(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) StripeCanceler {
	return &StripeCancel{
		GqlreqSvc: gqlreqSvc,
		GqlsdkSvc: gqlsdkSvc,
	}
}

//////////////////////////////////////
// Method implementation
//////////////////////////////////////

// Init the subscription
func (s *StripeCancel) Cancel(ctx context.Context, input *model.CancelInput) (*model.CancelOutput, error) {

	logrus.Debug("retrieving subscription id on hasura")
	subscriptionId, err := s.getStripeSubscriptionID(ctx, input.IDAccount)
	if err != nil {
		return nil, err
	}

	logrus.Debug("cancelling stripe subscription")
	ser, err := s.cancelStripeSubscription(ctx, input.IDAccount, subscriptionId)
	if err != nil {
		return nil, err
	}

	logrus.Debug("building cancel response")
	return &model.CancelOutput{
		Status: string(ser.Status),
	}, nil

}

// Retrieving the stripe subscription id from hasura
func (s *StripeCancel) getStripeSubscriptionID(ctx context.Context, accountID string) (string, error) {
	result, err := s.GqlsdkSvc.GetStripeSubscription(ctx, accountID)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_ACCOUNT_ID, accountID).Error("not able to execute GetStripeSubscription")
		return "", errorx.InternalError.Wrap(err, "not able to execute GetStripeSubscription")
	}

	if len(result.SubscriptionStatus) != 1 {
		logrus.WithContext(ctx).WithField(model.LOG_PARAM_RESULT_LENGTH, len(result.SubscriptionStatus)).Error("hasura GetStripeSubscription has returned wrong amount of results")
		return "", errorx.InternalError.Wrap(err, "subscription not able to cancelled, contact us")
	}

	logrus.WithContext(ctx).WithField(model.LOG_PARAM_HASURA_RESPONSE, logger.PrintStruct(result)).Debug("GetStripeSubscription result")

	return *result.SubscriptionStatus[0].StripeSubscriptionID, nil
}

// Requesting stripe to cancel the subscription
func (s *StripeCancel) cancelStripeSubscription(ctx context.Context, accountID string, stripeSubscriptioId string) (*stripe.Subscription, error) {
	subscriptionCancel, err := StripeCancelSubFunc(stripeSubscriptioId, nil)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			model.LOG_PARAM_ACCOUNT_ID:      accountID,
			model.LOG_PARAM_SUBSCRIPTION_ID: stripeSubscriptioId,
		}).Error("unable to cancel subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to process cancel request with the provider")
	}
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(subscriptionCancel),
		model.LOG_PARAM_ACCOUNT_ID:      accountID,
		model.LOG_PARAM_SUBSCRIPTION_ID: stripeSubscriptioId,
	}).Debug("subscription done")

	return subscriptionCancel, nil
}
