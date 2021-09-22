package saas

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

type getCurrentAccountHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewGetCurrentAccountHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &getCurrentAccountHandler{
		SdkSvc:     sdkSvc,
		GraphqlSvc: graphqlSvc,
	}
}

type ActionPayloadGetCurrentAccount struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.SaasGetCurrentAccountInput `json:"data"`
	} `json:"input"`
}

/*
Handle set current account for a user
*/
func (h *getCurrentAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("parsing payload")
	actionPayload := &ActionPayloadGetCurrentAccount{}
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

	setR, err := h.SdkSvc.GetCurrentAccount(r.Context(), authzInfo.UserId)
	if err != nil {
		const message = "not able to get the account account"
		logrus.WithError(err).WithFields(logrus.Fields{
			LOG_PARAM_USER_ID: authzInfo.UserId,
		}).Error(message)
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, message))
	}
	if len(setR.SaasMembership) > 1 {
		const message = "GetCurrentAccount returned wrong amount of data"
		logrus.WithError(err).WithFields(logrus.Fields{
			LOG_PARAM_USER_ID: authzInfo.UserId,
		}).Error(message)
		hshttp.WriteError(w, errorx.InternalError.New(message))
	}

	result := &gqlsdk.SaasGetCurrentAccountOutput{}

	if len(setR.SaasMembership) == 0 {
		logrus.Debug("building response for loggedin user")
		result.IDAccount = authz.ConfigAnonymousAccount()
		result.IDRole = authz.ConfigLoggedInRole()
	} else {
		logrus.Debug("building response from hasura")
		result.IDAccount = setR.SaasMembership[0].IDAccount
		result.IDRole = setR.SaasMembership[0].IDRole
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}

	logrus.WithFields(logrus.Fields{
		LOG_PARAM_ACCOUNT_ID: result.IDAccount,
		LOG_PARAM_USER_ID:    authzInfo.UserId,
		LOG_PARAM_ROLE_ID:    result.IDRole,
	}).Info("subscription init done")
}
