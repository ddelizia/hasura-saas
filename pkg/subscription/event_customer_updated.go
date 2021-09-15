package subscription

import (
	"github.com/stripe/stripe-go"
)

func ProcessCustomerUpdated(event stripe.Event, id string) error {
	data := &stripe.Customer{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
