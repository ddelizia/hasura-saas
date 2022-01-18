package subscription

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

// Get Stripe plan id from existing account in hasura, created during the init step
func getStripePlanFromPlan(ctx context.Context, sdkSvc gqlsdk.Service, id string) (*gqlsdk.QueryGetStripePlanFromPlan, error) {

	result, err := sdkSvc.GetStripePlanFromPlan(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_PLAN_ID, id).Error("stripe plan not found")
		return nil, errorx.InternalError.Wrap(err, "not able to get provider plan information")
	}

	if len(result.SubscriptionPlan) != 1 {
		logrus.WithFields(logrus.Fields{
			LOG_PARAM_HASURA_RESPONSE: len(result.SubscriptionPlan),
			LOG_PARAM_STRIPE_PLAN_ID:  id,
		}).Error("hasura GetStripePlanFromPlan has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "provider plan is not available")
	}

	logrus.WithField(LOG_PARAM_HASURA_RESPONSE, logger.PrintStruct(result)).Debug("GetStripePlanFromPlanID result")

	return result, nil
}

// Get Stripe customer id from existing account in hasura, created during the init step
func getStripeCustomerIdFromHasura(ctx context.Context, sdkSvc gqlsdk.Service, authz *gqlreq.HeaderInfo) (*gqlsdk.QueryGetAccountInfoForCreatingSubscription, error) {

	result, err := sdkSvc.GetAccountInfoForCreatingSubscription(ctx, authz.AccountId)
	if err != nil {
		logrus.WithError(err).Error("not able to get customer information from hasure")
		return nil, errorx.InternalError.Wrap(err, "not able to get customer information")
	}

	if len(result.SaasAccount) != 1 {
		logrus.WithField("wrong.result.length", len(result.SaasAccount)).Error("hasura QuerySubscription has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "customer information is wrong, contact us")
	}

	logrus.WithField("result", logger.PrintStruct(result)).Debug("QuerySubscription result")

	return result, nil
}

// Get Stripe plan id from existing account in hasura, created during the init step
func getPlanFromStripePlan(ctx context.Context, sdkSvc gqlsdk.Service, stripe_plan_id string) (*gqlsdk.QueryGetPlanFromStripePlan, error) {

	result, err := sdkSvc.GetPlanFromStripePlan(ctx, stripe_plan_id)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_STRIPE_PLAN_ID, stripe_plan_id).Error("plan not found")
		return nil, errorx.InternalError.Wrap(err, "not able to get plan information")
	}

	if len(result.SubscriptionPlan) != 1 {
		logrus.WithFields(logrus.Fields{
			LOG_PARAM_HASURA_RESPONSE: len(result.SubscriptionPlan),
			LOG_PARAM_STRIPE_PLAN_ID:  stripe_plan_id,
		}).Error("hasura GetPlanFromStripePlan has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "plan is not available")
	}

	logrus.WithField(LOG_PARAM_HASURA_RESPONSE, logger.PrintStruct(result)).Debug("GetPlanFromStripePlan result")

	return result, nil
}

// Updating Hasura subscription status information
func updateHasuraSubscription(ctx context.Context, sdkSvc gqlsdk.Service, accountId string, ser *stripe.Subscription) (*gqlsdk.MutationSetSubscriptioStatus, error) {

	plan, err := getPlanFromStripePlan(ctx, sdkSvc, ser.Plan.ID)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "unable to find plan while SetSubscriptioStatus")
	}

	result, err := sdkSvc.SetSubscriptioStatus(ctx, string(ser.Status), string(ser.Status) == "active" || string(ser.Status) == "trialing", accountId, ser.ID, plan.SubscriptionPlan[0].ID)
	if err != nil {
		logrus.WithError(err).Error("not able to execute mutation SetSubscriptioStatus")
		return nil, errorx.InternalError.Wrap(err, "unable to execute mutation SetSubscriptioStatus")
	}

	if len(result.UpdateSubscriptionStatus.Returning) != 1 {
		logrus.WithField("wrong.result.length", len(result.UpdateSubscriptionStatus.Returning)).Error("hasura SetSubscriptioStatus has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "subscription not able to be updated, contact us")
	}

	logrus.WithField("result", logger.PrintStruct(result)).Debug("SetSubscriptioStatus result")

	return result, nil
}
