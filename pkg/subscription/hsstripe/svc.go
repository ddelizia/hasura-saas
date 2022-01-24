package hsstripe

import (
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
)

type Service interface {
	StripeInitter
	StripeCreator
	StripeRetryer
}

type service struct {
	*StripeInit
	*StripeCreate
	*StripeRetry
}

func NewService(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) Service {
	stripeInit := NewStripeInit(gqlreqSvc, gqlsdkSvc).(*StripeInit)
	stripeCreate := NewStripeCreate(gqlreqSvc, gqlsdkSvc).(*StripeCreate)
	stripeRetry := NewStripeRetry(gqlreqSvc, gqlsdkSvc).(*StripeRetry)
	return &service{
		StripeInit:   stripeInit,
		StripeCreate: stripeCreate,
		StripeRetry:  stripeRetry,
	}
}
