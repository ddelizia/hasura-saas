package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/authz"
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
)

//////////////////////////////////////
// Interface
//////////////////////////////////////

type StripeInitter interface {
	Init(ctx context.Context, input *model.InitInput) (*model.InitOutput, error)
}

//////////////////////////////////////
// Struct
//////////////////////////////////////

type StripeInit struct {
	GqlreqSvc gqlreq.Service
	GqlsdkSvc gqlsdk.Service
	// wrap stripe function call for testing purpuses
	StripeNewCustomerFunc func(accountName string) (*stripe.Customer, error)
}

//////////////////////////////////////
// Mock
//////////////////////////////////////

type StripeInitMock struct {
	mock.Mock
}

func (m *StripeInitMock) Init(ctx context.Context, input *model.InitInput) (*model.InitOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.InitOutput), args.Error(1)
}

//////////////////////////////////////
// New
//////////////////////////////////////

func NewStripeInit(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) StripeInitter {
	return &StripeInit{
		GqlreqSvc:             gqlreqSvc,
		GqlsdkSvc:             gqlsdkSvc,
		StripeNewCustomerFunc: stripeNewCustomerFunc,
	}
}

//////////////////////////////////////
// Method implementation
//////////////////////////////////////

// Init the subscription
func (s *StripeInit) Init(ctx context.Context, input *model.InitInput) (*model.InitOutput, error) {

	logrus.WithContext(ctx).Debug("make sure plan exists on hasura")
	err := s.checkStripePlanRegisteredOnHasura(ctx, input.IDPlan)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("creating customer on stripe")
	c, err := s.createCustomerOnStripe(ctx, input.AccountName)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("updating user information on hasura")
	accountMutationResp, err := s.createCustomerSubscriptionOnHasura(ctx, input, c)
	if err != nil {
		return nil, err
	}

	logrus.WithContext(ctx).Debug("building init response")
	return &model.InitOutput{
		IDAccount: accountMutationResp.InsertSaasAccount.Returning[0].ID,
	}, nil
}

// Utility methods

func (s *StripeInit) checkStripePlanRegisteredOnHasura(ctx context.Context, id string) error {

	result, err := s.GqlsdkSvc.GetStripePlanFromPlan(ctx, id)
	if err != nil {
		return hserrorx.Wrap(
			err, errorx.InternalError,
			hserrorx.Fields{model.LOG_PARAM_PLAN_ID: id},
			"not able to get provider plan information", hstype.NewString("stripe plan not found"),
		)
	}

	return hserrorx.AssertTrue(len(result.SubscriptionPlan) == 1,
		hserrorx.Fields{
			model.LOG_PARAM_STRIPE_PLAN_ID:  id,
			model.LOG_PARAM_HASURA_RESPONSE: len(result.SubscriptionPlan),
		},
		"provider plan is not available", hstype.NewString("hasura GetStripePlanFromPlan has returned wrong amount of results"),
	)
}

func (s *StripeInit) createCustomerOnStripe(ctx context.Context, accountName string) (*stripe.Customer, error) {
	c, err := s.StripeNewCustomerFunc(accountName)

	if err != nil {
		return nil, hserrorx.Wrap(
			err, errorx.InternalError,
			hserrorx.Fields{model.LOG_PARAM_CUSTOMER_ID: c.ID},
			"unable to create customer", hstype.NewString("unable to create stripe customer due to an error on strype"),
		)
	}

	logrus.WithContext(ctx).WithFields(logrus.Fields{
		model.LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(c),
		model.LOG_PARAM_CUSTOMER_ID:     c.ID,
		model.LOG_PARAM_ACCOUNT_NAME:    accountName,
	}).Info("customer created on stripe")
	return c, nil
}

func (s *StripeInit) createCustomerSubscriptionOnHasura(ctx context.Context, i *model.InitInput, c *stripe.Customer) (*gqlsdk.MutationCreateSubscriptionCustomer, error) {

	customer, err := s.GqlsdkSvc.CreateSubscriptionCustomer(
		ctx,
		i.AccountName,
		i.IDPlan,
		i.IDUser,
		c.ID,
		model.STATUS_INIT,
		authz.ConfigAccountOwnerRole(),
	)

	if err != nil {
		return nil, hserrorx.Wrap(
			err, errorx.InternalError,
			hserrorx.Fields{model.LOG_PARAM_USER_ID: i.IDUser},
			"unable to create customer", hstype.NewString("unable to create customer on hasura"),
		)
	}
	return customer, nil
}

func stripeNewCustomerFunc(accountName string) (*stripe.Customer, error) {
	params := stripe.CustomerParams{
		Description: stripe.String("Stripe customer for account " + accountName),
	}
	return customer.New(&params)
}
