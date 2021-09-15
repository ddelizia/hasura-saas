package subscription

import (
	"github.com/stripe/stripe-go"
)

func ProcessInvoicePaid(event stripe.Event, id string) error {
	data := &stripe.Invoice{}
	if err := beforeEvent(event, data); err != nil {
		return err
	}

	return nil
}
