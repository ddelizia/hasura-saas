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

type ChangeHandler struct {
	StripeChanger hsstripe.StripeChanger
}

func NewChangeHandler(stripeChanger hsstripe.StripeChanger) http.Handler {
	return &ChangeHandler{
		StripeChanger: stripeChanger,
	}
}

type ActionChangePayload struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.ChangeSubscriptionInput `json:"data"`
	} `json:"input"`
}

/*
Handle subscription change
*/
func (h *ChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := hscontext.ActionDataValue(r.Context()).(*gqlsdk.ChangeSubscriptionInput)
	authzInfo := hscontext.AuthzInfoValue(r.Context())

	got, err := h.StripeChanger.Change(r.Context(), &model.ChangeInput{
		IDAccount: authzInfo.AccountId,
		IDPlan:    data.IDPlan,
	})

	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to execute stripe change"))
		return
	}

	out := &gqlsdk.ChangeSubscriptionOutput{
		IDAccount: got.IDAccount,
		IsActive:  got.IsActive,
	}

	err = hshttp.WriteBody(w, out)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}

}
