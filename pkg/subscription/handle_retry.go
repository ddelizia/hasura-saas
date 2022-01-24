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

type RetryHandler struct {
	StripeRetryer hsstripe.StripeRetryer
}

func NewRetryHandler(retryer hsstripe.StripeRetryer) http.Handler {
	return &RetryHandler{
		StripeRetryer: retryer,
	}
}

type ActionPayloadRetry struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.RetrySubscriptionInput `json:"data"`
	} `json:"input"`
}

func (h *RetryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := hscontext.ActionDataValue(r.Context()).(*gqlsdk.RetrySubscriptionInput)
	authzInfo := hscontext.AuthzInfoValue(r.Context())

	got, err := h.StripeRetryer.Retry(r.Context(), &model.RetryInput{
		IDAccount:       authzInfo.AccountId,
		IDPaymentMethod: &data.PaymentMethodID,
	})

	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to execute stripe retry"))
		return
	}

	out := &gqlsdk.RetrySubscriptionOutput{
		IDAccount: got.IDAccount,
		IsActive:  got.IsActive,
	}

	err = hshttp.WriteBody(w, out)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}
