package subscription

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

/*
 Ger Stripe plan id from existing account in hasura, created during the init step
*/
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
