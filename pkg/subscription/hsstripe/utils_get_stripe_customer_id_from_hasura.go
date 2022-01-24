package hsstripe

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

// Get Stripe customer id from existing account in hasura, created during the init step
func getStripeCustomerIdFromHasura(ctx context.Context, sdkSvc gqlsdk.Service, accountId string) (*gqlsdk.QueryGetAccountInfoForCreatingSubscription, error) {

	result, err := sdkSvc.GetAccountInfoForCreatingSubscription(ctx, accountId)
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
