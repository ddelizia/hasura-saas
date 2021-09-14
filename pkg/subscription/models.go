package subscription

import "github.com/ddelizia/hasura-saas/pkg/hstype"

type RequestCreate struct {
	PaymentMethodID string `json:"paymentMethodId"`
	CustomerID      string `json:"customerId"`
	PriceID         string `json:"priceId"`
	AccountName     string `json:"accountName"`
}

type RequestRetry struct {
	CustomerID      string `json:"customerId"`
	PaymentMethodID string `json:"paymentMethodId"`
	InvoiceID       string `json:"invoiceId"`
	AccountID       string `json:"accountId"`
}

type RequestCancel struct {
	SubscriptionID string `json:"subscriptionId"`
}

type Plan struct {
	ID hstype.String
}
