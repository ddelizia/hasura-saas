package subscription

import (
	"github.com/stripe/stripe-go"
)

func ProcessPaymentMethodAttached(event stripe.Event, id string) error {
	data := &stripe.PaymentMethod{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
