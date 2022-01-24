package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go"
)

//////////////////////////////////////
// Interface
//////////////////////////////////////

type StripeChanger interface {
	Change(ctx context.Context, input *model.ChangeInput) (*model.ChangeOutput, error)
}

//////////////////////////////////////
// Struct
//////////////////////////////////////

type StripeChange struct {
	GqlreqSvc gqlreq.Service
	GqlsdkSvc gqlsdk.Service
}

//////////////////////////////////////
// Mock
//////////////////////////////////////

type StripeChangeMock struct {
	mock.Mock
}

func (m *StripeChangeMock) Change(ctx context.Context, input *model.ChangeInput) (*model.ChangeOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.ChangeOutput), args.Error(1)
}

//////////////////////////////////////
// New
//////////////////////////////////////

func NewStripeChange(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) StripeChanger {
	return &StripeChange{
		GqlreqSvc: gqlreqSvc,
		GqlsdkSvc: gqlsdkSvc,
	}
}

//////////////////////////////////////
// Method implementation
//////////////////////////////////////

// Init the subscription
func (s *StripeChange) Change(ctx context.Context, input *model.ChangeInput) (*model.ChangeOutput, error) {

	logrus.Debug("getting saas account information from hasura")
	accountInfoForCreatingSubscription, err := getStripeCustomerIdFromHasura(ctx, s.GqlsdkSvc, input.IDAccount)
	if err != nil {
		return nil, err
	}

	logrus.Debug("get stripe plan id")
	plan, err := getStripePlanFromPlan(ctx, s.GqlsdkSvc, input.IDPlan)
	if err != nil {
		return nil, err
	}

	logrus.Debug("update subscription on stripe")
	ser, err := s.updateStripeSubscription(
		ctx,
		*accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionStatus.StripeSubscriptionID,
		*plan.SubscriptionPlan[0].StripeCode)
	if err != nil {
		return nil, err
	}

	logrus.Debug("create subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(ctx, s.GqlsdkSvc, input.IDAccount, ser)
	if err != nil {
		return nil, err
	}

	logrus.Debug("building response")
	return &model.ChangeOutput{
		IDAccount: input.IDAccount,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}, nil

}

/*
Update subscription on stripe
* Get subscription information
* Update subscription plan
*/
func (s *StripeChange) updateStripeSubscription(ctx context.Context, subscriptionId string, planId string) (*stripe.Subscription, error) {
	logrus.Debug("getting subscription information from stripe")
	sub, err := StripeGetSubFunc(subscriptionId, &stripe.SubscriptionParams{})
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_SUBSCRIPTION_ID, subscriptionId).Error("payment attachment failde")
		return nil, errorx.InternalError.Wrap(err, "unable to retrieve subscription")
	}

	logrus.Debug("updating subscription on stripe")
	updatedSubscription, err := StripeUpdateSubFunc(
		subscriptionId,
		&stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(false),
			Items: []*stripe.SubscriptionItemsParams{{
				ID:   stripe.String(sub.Items.Data[0].ID),
				Plan: stripe.String(planId),
			}},
		},
	)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_SUBSCRIPTION_ID, subscriptionId).Error("not able to update subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to update subscription")
	}

	return updatedSubscription, nil
}
