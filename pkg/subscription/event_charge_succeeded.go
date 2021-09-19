package subscription

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/stripe/stripe-go"
)

func ProcessChargeSucceeded(c context.Context, event stripe.Event, id string, sdkSvc gqlsdk.Service) error {
	data := &stripe.Charge{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
