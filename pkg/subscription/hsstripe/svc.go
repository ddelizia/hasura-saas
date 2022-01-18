package hsstripe

import (
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
)

type Service interface {
	StripeInitter
	StripeCreator
}

type service struct {
	*StripeInit
	*StripeCreate
}

func NewService(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) Service {
	stripeInit := NewStripeInit(gqlreqSvc, gqlsdkSvc).(*StripeInit)
	stripeCreate := NewStripeCreate(gqlreqSvc, gqlsdkSvc).(*StripeCreate)
	return &service{
		StripeInit:   stripeInit,
		StripeCreate: stripeCreate,
	}
}
