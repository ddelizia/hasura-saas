package subscription

import (
	"github.com/stripe/stripe-go"
)

func ProcessPaymentIntentRequiresAction(event stripe.Event, id string) error {
	data := &stripe.PaymentIntent{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
