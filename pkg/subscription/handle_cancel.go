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

type CancelHandler struct {
	StripeCanceler hsstripe.StripeCanceler
}

func NewCancelHandler(canceler hsstripe.StripeCanceler) http.Handler {
	return &CancelHandler{
		StripeCanceler: canceler,
	}
}

type ActionPayloadCancel struct {
	hasura.BasePayload
}

func (h *CancelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	authzInfo := hscontext.AuthzInfoValue(r.Context())

	got, err := h.StripeCanceler.Cancel(r.Context(), &model.CancelInput{
		IDAccount: authzInfo.AccountId,
	})

	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to execute stripe cancel"))
		return
	}

	out := &gqlsdk.CancelSubscriptionOutput{
		Status: got.Status,
	}

	err = hshttp.WriteBody(w, out)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}
