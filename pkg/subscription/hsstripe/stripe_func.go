package hsstripe

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/paymentmethod"
	"github.com/stripe/stripe-go/sub"
)

// Wrapping stripe functions for testing purposes
var (
	// Customer

	StripeUpdateCustomerFunc func(id string, params *stripe.CustomerParams) (*stripe.Customer, error) = customer.Update
	StripeNewCustomerFunc    func(params *stripe.CustomerParams) (*stripe.Customer, error)            = customer.New

	// Invoice

	StripeListInvoiceFunc func(params *stripe.InvoiceListParams) *invoice.Iter                      = invoice.List
	StripeGetInvoiceFunc  func(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error)    = invoice.Get
	StripePayInvoiceFunc  func(id string, params *stripe.InvoicePayParams) (*stripe.Invoice, error) = invoice.Pay

	// Payment

	StripeAttachPaymentmethodFunc func(id string, params *stripe.PaymentMethodAttachParams) (*stripe.PaymentMethod, error) = paymentmethod.Attach

	// Subscription

	StripeNewSubFunc    func(params *stripe.SubscriptionParams) (*stripe.Subscription, error)                  = sub.New
	StripeGetSubFunc    func(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error)       = sub.Get
	StripeCancelSubFunc func(id string, params *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) = sub.Cancel
)
