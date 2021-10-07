package subscription

import (
	"net/http"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/gqlsdk"
	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/sub"
)

type ChangeHandler struct {
	GraphqlSvc gqlreq.Service
	SdkSvc     gqlsdk.Service
}

func NewChangeHandler(graphqlSvc gqlreq.Service, sdkSvc gqlsdk.Service) http.Handler {
	return &ChangeHandler{
		GraphqlSvc: graphqlSvc,
		SdkSvc:     sdkSvc,
	}
}

type ActionChangePayload struct {
	hasura.BasePayload
	Input struct {
		Data *gqlsdk.ChangeSubscriptionInput `json:"data"`
	} `json:"input"`
}

// channels
type subcriptionChanType struct {
	subscriptionId string
	authzInfo      *gqlreq.HeaderInfo
	err            error
}

type planChanType struct {
	planId string
	err    error
}

/*
Handle subscription change
*/
func (h *ChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logrus.Debug("change subscription request")
	actionPayload := &ActionChangePayload{}
	err := hshttp.GetBody(r, actionPayload)
	if err != nil {
		hshttp.WriteError(w, errorx.IllegalArgument.Wrap(err, "invalid payload for change subscription"))
		return
	}

	subscriptionChan := make(chan subcriptionChanType)
	planChan := make(chan planChanType)

	go func() {
		defer close(subscriptionChan)

		logrus.Debug("getting authorization information")
		authzInfo, err := h.GraphqlSvc.GetSessionInfo(actionPayload.SessionVariables)
		if err != nil {
			logrus.Error("unable to retrieve authz information")
			hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve authz information"))
			subscriptionChan <- subcriptionChanType{err: err}
			return
		}

		logrus.Debug("getting saas account information from hasura")
		accountInfoForCreatingSubscription, err := getStripeCustomerIdFromHasura(r.Context(), h.SdkSvc, authzInfo)
		if err != nil {
			logrus.Error("unable to retrieve customer")
			hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to retrieve customer"))
			subscriptionChan <- subcriptionChanType{err: err}
			return
		}

		subscriptionChan <- subcriptionChanType{
			subscriptionId: *accountInfoForCreatingSubscription.SaasAccount[0].SubscriptionStatus.StripeSubscriptionID,
			authzInfo:      authzInfo,
			err:            nil,
		}
	}()

	go func() {
		defer close(planChan)

		logrus.Debug("get stripe plan id")
		plan, err := getStripePlanFromPlan(r.Context(), h.SdkSvc, actionPayload.Input.Data.IDPlan)
		if err != nil {
			logrus.Error("plan not found")
			hshttp.WriteError(w, errorx.InternalError.Wrap(err, "plan not found"))
			planChan <- planChanType{err: err}
			return
		}

		planChan <- planChanType{
			planId: *plan.SubscriptionPlan[0].StripeCode,
			err:    nil,
		}
	}()

	subscription, plan := <-subscriptionChan, <-planChan

	if subscription.err != nil || plan.err != nil {
		return
	}

	logrus.Debug("update subscription on stripe")
	ser, err := h.updateStripeSubscription(subscription.subscriptionId, plan.planId)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to update information"))
		return
	}

	logrus.Debug("create subscription on hasura")
	updatedStatus, err := updateHasuraSubscription(r.Context(), h.SdkSvc, subscription.authzInfo.AccountId, ser)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "unable to store subscription information"))
		return
	}

	logrus.Debug("building response")
	result := &gqlsdk.ChangeSubscriptionOutput{
		IDAccount: subscription.authzInfo.AccountId,
		IsActive:  updatedStatus.UpdateSubscriptionStatus.Returning[0].IsActive,
	}

	err = hshttp.WriteBody(w, result)
	if err != nil {
		hshttp.WriteError(w, errorx.InternalError.Wrap(err, "not able to create response"))
		return
	}

	logrus.WithFields(logrus.Fields{
		LOG_PARAM_ACCOUNT_ID:      subscription.authzInfo.AccountId,
		LOG_PARAM_SUBSCRIPTION_ID: ser.ID,
		LOG_PARAM_CUSTOMER_ID:     ser.Customer.ID,
		LOG_PARAM_PLAN_ID:         ser.Plan.ID,
	}).Info("subscription change done")
}

/*
Update subscription on stripe
* Get subscription information
* Update subscription plan
*/
func (h *ChangeHandler) updateStripeSubscription(subscriptionId string, planId string) (*stripe.Subscription, error) {
	logrus.Debug("getting subscription information from stripe")
	s, err := sub.Get(subscriptionId, nil)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_SUBSCRIPTION_ID, subscriptionId).Error("payment attachment failde")
		return nil, errorx.InternalError.Wrap(err, "unable to retrieve subscription")
	}

	logrus.Debug("updating subscription on stripe")
	updatedSubscription, err := sub.Update(
		subscriptionId,
		&stripe.SubscriptionParams{
			CancelAtPeriodEnd: stripe.Bool(false),
			Items: []*stripe.SubscriptionItemsParams{{
				ID:   stripe.String(s.Items.Data[0].ID),
				Plan: stripe.String(planId),
			}},
		},
	)
	if err != nil {
		logrus.WithError(err).WithField(LOG_PARAM_SUBSCRIPTION_ID, subscriptionId).Error("not able to update subscription")
		return nil, errorx.InternalError.Wrap(err, "unable to update subscription")
	}

	return updatedSubscription, nil
}
