package subscription

import (
	"context"
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/authz"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/ddelizia/hasura-saas/pkg/logger"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

type initHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewInitHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &initHandler{
		SdkSvc:     sdkSvc,
		GraphqlSvc: graphqlSvc,
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

	logrus.Debug("parsing payload")
	actionPayload := &ActionPayloadInit{}
	err := hshttp.GetBody(r, actionPayload)
	if err != nil {
		hshttp.WriteError(w, errorx.IllegalArgument.Wrap(err, "invalid payload for create customer"))
		return
	}

	logrus.Debug("get authorization info from session")
	authzInfo, err := h.GraphqlSvc.GetSessionInfo(actionPayload.SessionVariables)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve authz information"))
		return
	}

	_, err = getStripePlanFromPlan(r.Context(), h.SdkSvc, actionPayload.Input.Data.IDPlan)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "plan not found"))
		return
	}

	logrus.Debug("creating stripe customer")
	c, err := h.createStripeCustomer(actionPayload.Input.Data.AccountName)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to create customer"))
		return
	}

	logrus.Debug("updating user information on hasura")
	accountMutationResp, err := h.createCustomerSubscriptionOnHasura(r.Context(), actionPayload, authzInfo, c)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to create account information"))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.InitSubscriptionOutput{
		AccountID: accountMutationResp.InsertSaasAccount.Returning[0].ID,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}
}

/*
 Create stripe customer
*/
func (s *initHandler) createStripeCustomer(accountName string) (*stripe.Customer, error) {
	params := stripe.CustomerParams{
		Description: stripe.String("Stripe customer for account " + accountName),
	}
	c, err := customer.New(&params)

	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_CUSTOMER_ID, c.ID).Error("unable to create customer")
		return nil, errorx.InternalError.Wrap(err, "unable to create customer")
	}

	logrus.WithFields(logrus.Fields{
		LOG_PARAM_STRIPE_RESPONSE: logger.PrintStruct(c),
		LOG_PARAM_CUSTOMER_ID:     c.ID,
		LOG_PARAM_ACCOUNT_NAME:    accountName,
	}).Info("stripe customer created")
	return c, nil

}

/*
 Create customer subscription on hasura
 At this point:
 * Account is being created
 * Calling user is being added as account admin
 * Stripe customer is being added to the account
*/
func (h *initHandler) createCustomerSubscriptionOnHasura(ctx context.Context, a *ActionPayloadInit, az *gqlreq.HeaderInfo, c *stripe.Customer) (*gqlsdk.MutationCreateSubscriptionCustomer, error) {

	customer, err := h.SdkSvc.CreateSubscriptionCustomer(
		ctx,
		a.Input.Data.AccountName,
		a.Input.Data.IDPlan,
		az.UserId,
		c.ID,
		STATUS_INIT,
		authz.ConfigAccountOwnerRole(),
	)

	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_USER_ID, az.UserId).Error("unable to create customer on hasura")
		return nil, errorx.InternalError.Wrap(err, "unable to execute CreateSubscriptionCustomer")
	}
	return customer, nil
}
