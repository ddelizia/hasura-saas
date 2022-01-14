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

type initHandler struct {
	StripeInitter hsstripe.StripeInitter
}

func NewInitHandler(initter hsstripe.StripeInitter) http.Handler {
	return &initHandler{
		StripeInitter: initter,
	}
}

type ActionPayloadInit struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.InitSubscriptionInput `json:"data"`
	} `json:"input"`
}

/*
Handle subscription initialization
*/
func (h *initHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := hscontext.ActionDataValue(r.Context()).(*ActionPayloadInit)
	authzInfo := hscontext.AuthzInfoValue(r.Context())

	out, err := h.StripeInitter.Init(r.Context(), &model.InitInput{
		AccountName: data.Input.Data.AccountName,
		IDPlan:      data.Input.Data.IDPlan,
		IDUser:      authzInfo.UserId,
	})

	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to execute stripe initialization"))
		return
	}

	err = hshttp.WriteBody(w, out)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}
