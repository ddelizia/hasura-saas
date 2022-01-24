package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

// Updating Hasura subscription status information
func updateHasuraSubscription(ctx context.Context, sdkSvc gqlsdk.Service, accountId string, ser *stripe.Subscription) (*gqlsdk.MutationSetSubscriptioStatus, error) {

	plan, err := getPlanFromStripePlan(ctx, sdkSvc, ser.Plan.ID)
	if err != nil {
		return nil, errorx.InternalError.Wrap(err, "unable to find plan while SetSubscriptioStatus")
	}

	result, err := sdkSvc.SetSubscriptioStatus(ctx, string(ser.Status), string(ser.Status) == "active" || string(ser.Status) == "trialing", accountId, ser.ID, plan.SubscriptionPlan[0].ID)
	if err != nil {
		logrus.WithContext(ctx).WithError(err).Error("not able to execute mutation SetSubscriptioStatus")
		return nil, errorx.InternalError.Wrap(err, "unable to execute mutation SetSubscriptioStatus")
	}

	if len(result.UpdateSubscriptionStatus.Returning) != 1 {
		logrus.WithContext(ctx).WithField("wrong.result.length", len(result.UpdateSubscriptionStatus.Returning)).Error("hasura SetSubscriptioStatus has returned wrong amount of results")
		return nil, errorx.InternalError.Wrap(err, "subscription not able to be updated, contact us")
	}

	logrus.WithContext(ctx).WithField("result", logger.PrintStruct(result)).Debug("SetSubscriptioStatus result")

	// add comment

	return result, nil
}
