package saas

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

type setCurrentAccountHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewSetCurrentAccountHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &setCurrentAccountHandler{
		SdkSvc:     sdkSvc,
		GraphqlSvc: graphqlSvc,
	}
}

type ActionPayloadSetCurrentAccount struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.SaasSetCurrentAccountInput `json:"data"`
	} `json:"input"`
}

/*
Handle set current account for a user
*/
func (h *setCurrentAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("parsing payload")
	actionPayload := &ActionPayloadSetCurrentAccount{}
	err := hshttp.GetBody(r, actionPayload)
	if err != nil {
		hshttp.WriteError(w, errorx.IllegalArgument.Wrap(err, "invalid payload for set current account"))
		return
	}

	logrus.Debug("get authorization info from session")
	authzInfo, err := h.GraphqlSvc.GetSessionInfo(actionPayload.SessionVariables)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve authz information"))
		return
	}

	setR, err := h.SdkSvc.SetAccountForUser(r.Context(), actionPayload.Input.Data.IDAccount, authzInfo.UserId)
	if err != nil {
		const message = "not able to set the account"
		logrus.WithError(err).WithFields(logrus.Fields{
			LOG_PARAM_ACCOUNT_ID: actionPayload.Input.Data.IDAccount,
			LOG_PARAM_USER_ID:    authzInfo.UserId,
		}).Error(message)
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, message))
		return
	}
	if len(setR.UpdateSaasMembership.Returning) != 1 {
		const message = "SetAccountForUser returned wrong amount of data"
		logrus.WithError(err).WithFields(logrus.Fields{
			LOG_PARAM_ACCOUNT_ID: actionPayload.Input.Data.IDAccount,
			LOG_PARAM_USER_ID:    authzInfo.UserId,
		}).Error(message)
		hshttp.WriteError(w, errorx.InternalError.New(message))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.SaasSetCurrentAccountOutput{
		IDAccount: actionPayload.Input.Data.IDAccount,
		//*result.UpdateSaasMembership.Returning[0].SelectedAt,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}

	logrus.WithFields(logrus.Fields{
		LOG_PARAM_ACCOUNT_ID: actionPayload.Input.Data.IDAccount,
		LOG_PARAM_USER_ID:    authzInfo.UserId,
	}).Info("set current account done")
}
