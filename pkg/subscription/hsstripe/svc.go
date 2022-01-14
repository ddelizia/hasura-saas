package hsstripe

import (
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
)

type Service interface {
	StripeInitter
}

type service struct {
	*StripeInit
}

func NewService(gqlreqSvc gqlreq.Service, gqlsdkSvc gqlsdk.Service) Service {
	stripeInit := NewStripeInit(gqlreqSvc, gqlsdkSvc).(*StripeInit)
	return &service{
		StripeInit: stripeInit,
	}
}
