package subscription

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/subscription/hsstripe"
	"github.com/ddelizia/hasura-saas/pkg/subscription/model"
	"github.com/joomcode/errorx"
)

type CreateHandler struct {
	Stripe hsstripe.StripeCreator
}

func NewCreateHandler(stripe hsstripe.StripeCreator) http.Handler {
	return &CreateHandler{
		Stripe: stripe,
	}
}

type ActionPayloadCreate struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.CreateSubscriptionInput `json:"data"`
	} `json:"input"`
}

/*
Handle subscription creation
*/
func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := hscontext.ActionDataValue(r.Context()).(*gqlsdk.CreateSubscriptionInput)
	authzInfo := hscontext.AuthzInfoValue(r.Context())

	got, err := h.Stripe.Create(r.Context(), &model.CreateInput{
		IDAccount:       authzInfo.AccountId,
		IDPaymentMethod: *data.PaymentMethodID,
	})

	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to execute stripe creation"))
		return
	}

	out := &gqlsdk.CreateSubscriptionOutput{
		IDAccount: got.IDAccount,
		IsActive:  got.IsActive,
	}
	err = hshttp.WriteBody(w, out)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}
