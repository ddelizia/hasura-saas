package subscription

import (
	"github.com/stripe/stripe-go"
)

func ProcessChargeSucceeded(event stripe.Event, id string) error {
	data := &stripe.Charge{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
