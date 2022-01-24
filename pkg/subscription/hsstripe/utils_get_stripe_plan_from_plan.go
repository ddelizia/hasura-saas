package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

func getStripePlanFromPlan(ctx context.Context, sdkSvc gqlsdk.Service, id string) (*gqlsdk.QueryGetStripePlanFromPlan, error) {

	result, err := sdkSvc.GetStripePlanFromPlan(ctx, id)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField(model.LOG_PARAM_PLAN_ID, id).Error("stripe plan not found")
		return nil, errorx.InternalError.Wrap(err, "not able to get provider plan information")
	}

	if len(result.SubscriptionPlan) != 1 {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			model.LOG_PARAM_HASURA_RESPONSE: len(result.SubscriptionPlan),
			model.LOG_PARAM_STRIPE_PLAN_ID:  id,
		}).Error("hasura GetStripePlanFromPlan has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "provider plan is not available")
	}

	logrus.WithContext(ctx).WithField(model.LOG_PARAM_HASURA_RESPONSE, logger.PrintStruct(result)).Debug("GetStripePlanFromPlanID result")

	return result, nil
}
