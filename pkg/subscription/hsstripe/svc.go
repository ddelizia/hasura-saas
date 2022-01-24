package hsstripe

import (
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
)

type Service interface {
	StripeInitter
	StripeCreator
	StripeRetryer
	StripeCanceler
	StripeChanger
}

type service struct {
	*StripeInit
	*StripeCreate
	*StripeRetry
	*StripeCancel
	*StripeChange
}

func NewService(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) Service {
	stripeInit := NewStripeInit(gqlreqSvc, gqlsdkSvc).(*StripeInit)
	stripeCreate := NewStripeCreate(gqlreqSvc, gqlsdkSvc).(*StripeCreate)
	stripeRetry := NewStripeRetry(gqlreqSvc, gqlsdkSvc).(*StripeRetry)
	stripeCancel := NewStripeCancel(gqlreqSvc, gqlsdkSvc).(*StripeCancel)
	stripeChange := NewStripeChange(gqlreqSvc, gqlsdkSvc).(*StripeChange)
	return &service{
		StripeInit:   stripeInit,
		StripeCreate: stripeCreate,
		StripeRetry:  stripeRetry,
		StripeCancel: stripeCancel,
		StripeChange: stripeChange,
	}
}
