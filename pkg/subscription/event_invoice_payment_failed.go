package subscription

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

func ProcessInvoicePaymentFailed(c context.Context, event stripe.Event, id string, sdkSvc gqlsdk.Service) error {
	data := &stripe.Invoice{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	logrus.WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).
		Debug("get account id for subscription")
	accountInfo, err := sdkSvc.GetAccountFromSubscription(c, data.Subscription.ID)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).Error("unable to get account id")
		return errorx.InternalError.Wrap(err, "unable to get account id")
	}
	if len(accountInfo.SubscriptionStatus) != 1 {
		logrus.WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).Error("multiple accounts have been returned")
		return errorx.InternalError.New("unable to get account id")
	}

	logrus.WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).
		Debug("get subscription information from invoice")
	subscriptionParams := &stripe.SubscriptionParams{}
	ser, err := sub.Get(data.Subscription.ID, subscriptionParams)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).Error("unable to get subscription")
		return errorx.InternalError.Wrap(err, "unable to get subscription")
	}

	logrus.WithField(LOG_PARAM_SUBSCRIPTION_ID, data.Subscription.ID).
		Debug("update status od the subscription on hasura")
	_, err = updateHasuraSubscription(c, sdkSvc, accountInfo.SubscriptionStatus[0].IDAccount, ser)
	if err != nil {
		return err
	}

	return nil
}
