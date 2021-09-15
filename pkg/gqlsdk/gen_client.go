// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package gqlsdk

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
)

type Client struct {
	Client *clientv2.Client
}

func NewClient(cli *http.Client, baseURL string, interceptors ...clientv2.RequestInterceptor) *Client {
	return &Client{Client: clientv2.NewClient(cli, baseURL, interceptors...)}
}

type QueryRoot struct {
	SaasAccount                     []*SaasAccount                  "json:\"saas_account\" graphql:\"saas_account\""
	SaasAccountAggregate            SaasAccountAggregate            "json:\"saas_account_aggregate\" graphql:\"saas_account_aggregate\""
	SaasAccountByPk                 *SaasAccount                    "json:\"saas_account_by_pk\" graphql:\"saas_account_by_pk\""
	SaasAddress                     []*SaasAddress                  "json:\"saas_address\" graphql:\"saas_address\""
	SaasAddressAggregate            SaasAddressAggregate            "json:\"saas_address_aggregate\" graphql:\"saas_address_aggregate\""
	SaasAddressByPk                 *SaasAddress                    "json:\"saas_address_by_pk\" graphql:\"saas_address_by_pk\""
	SaasMembership                  []*SaasMembership               "json:\"saas_membership\" graphql:\"saas_membership\""
	SaasMembershipAggregate         SaasMembershipAggregate         "json:\"saas_membership_aggregate\" graphql:\"saas_membership_aggregate\""
	SaasMembershipByPk              *SaasMembership                 "json:\"saas_membership_by_pk\" graphql:\"saas_membership_by_pk\""
	SaasRole                        []*SaasRole                     "json:\"saas_role\" graphql:\"saas_role\""
	SaasRoleAggregate               SaasRoleAggregate               "json:\"saas_role_aggregate\" graphql:\"saas_role_aggregate\""
	SaasRoleByPk                    *SaasRole                       "json:\"saas_role_by_pk\" graphql:\"saas_role_by_pk\""
	SaasUserAccount                 []*SaasUserAccount              "json:\"saas_user_account\" graphql:\"saas_user_account\""
	SaasUserAccountAggregate        SaasUserAccountAggregate        "json:\"saas_user_account_aggregate\" graphql:\"saas_user_account_aggregate\""
	SubscriptionActivePlan          []*SubscriptionActivePlan       "json:\"subscription_active_plan\" graphql:\"subscription_active_plan\""
	SubscriptionActivePlanAggregate SubscriptionActivePlanAggregate "json:\"subscription_active_plan_aggregate\" graphql:\"subscription_active_plan_aggregate\""
	SubscriptionCustomer            []*SubscriptionCustomer         "json:\"subscription_customer\" graphql:\"subscription_customer\""
	SubscriptionCustomerAggregate   SubscriptionCustomerAggregate   "json:\"subscription_customer_aggregate\" graphql:\"subscription_customer_aggregate\""
	SubscriptionCustomerByPk        *SubscriptionCustomer           "json:\"subscription_customer_by_pk\" graphql:\"subscription_customer_by_pk\""
	SubscriptionEvent               []*SubscriptionEvent            "json:\"subscription_event\" graphql:\"subscription_event\""
	SubscriptionEventAggregate      SubscriptionEventAggregate      "json:\"subscription_event_aggregate\" graphql:\"subscription_event_aggregate\""
	SubscriptionEventByPk           *SubscriptionEvent              "json:\"subscription_event_by_pk\" graphql:\"subscription_event_by_pk\""
	SubscriptionPlan                []*SubscriptionPlan             "json:\"subscription_plan\" graphql:\"subscription_plan\""
	SubscriptionPlanAggregate       SubscriptionPlanAggregate       "json:\"subscription_plan_aggregate\" graphql:\"subscription_plan_aggregate\""
	SubscriptionPlanByPk            *SubscriptionPlan               "json:\"subscription_plan_by_pk\" graphql:\"subscription_plan_by_pk\""
	SubscriptionStatus              []*SubscriptionStatus           "json:\"subscription_status\" graphql:\"subscription_status\""
	SubscriptionStatusAggregate     SubscriptionStatusAggregate     "json:\"subscription_status_aggregate\" graphql:\"subscription_status_aggregate\""
	SubscriptionStatusByPk          *SubscriptionStatus             "json:\"subscription_status_by_pk\" graphql:\"subscription_status_by_pk\""
}
type MutationRoot struct {
	DeleteSaasAccount               *SaasAccountMutationResponse            "json:\"delete_saas_account\" graphql:\"delete_saas_account\""
	DeleteSaasAccountByPk           *SaasAccount                            "json:\"delete_saas_account_by_pk\" graphql:\"delete_saas_account_by_pk\""
	DeleteSaasAddress               *SaasAddressMutationResponse            "json:\"delete_saas_address\" graphql:\"delete_saas_address\""
	DeleteSaasAddressByPk           *SaasAddress                            "json:\"delete_saas_address_by_pk\" graphql:\"delete_saas_address_by_pk\""
	DeleteSaasMembership            *SaasMembershipMutationResponse         "json:\"delete_saas_membership\" graphql:\"delete_saas_membership\""
	DeleteSaasMembershipByPk        *SaasMembership                         "json:\"delete_saas_membership_by_pk\" graphql:\"delete_saas_membership_by_pk\""
	DeleteSaasRole                  *SaasRoleMutationResponse               "json:\"delete_saas_role\" graphql:\"delete_saas_role\""
	DeleteSaasRoleByPk              *SaasRole                               "json:\"delete_saas_role_by_pk\" graphql:\"delete_saas_role_by_pk\""
	DeleteSubscriptionActivePlan    *SubscriptionActivePlanMutationResponse "json:\"delete_subscription_active_plan\" graphql:\"delete_subscription_active_plan\""
	DeleteSubscriptionCustomer      *SubscriptionCustomerMutationResponse   "json:\"delete_subscription_customer\" graphql:\"delete_subscription_customer\""
	DeleteSubscriptionCustomerByPk  *SubscriptionCustomer                   "json:\"delete_subscription_customer_by_pk\" graphql:\"delete_subscription_customer_by_pk\""
	DeleteSubscriptionEvent         *SubscriptionEventMutationResponse      "json:\"delete_subscription_event\" graphql:\"delete_subscription_event\""
	DeleteSubscriptionEventByPk     *SubscriptionEvent                      "json:\"delete_subscription_event_by_pk\" graphql:\"delete_subscription_event_by_pk\""
	DeleteSubscriptionPlan          *SubscriptionPlanMutationResponse       "json:\"delete_subscription_plan\" graphql:\"delete_subscription_plan\""
	DeleteSubscriptionPlanByPk      *SubscriptionPlan                       "json:\"delete_subscription_plan_by_pk\" graphql:\"delete_subscription_plan_by_pk\""
	DeleteSubscriptionStatus        *SubscriptionStatusMutationResponse     "json:\"delete_subscription_status\" graphql:\"delete_subscription_status\""
	DeleteSubscriptionStatusByPk    *SubscriptionStatus                     "json:\"delete_subscription_status_by_pk\" graphql:\"delete_subscription_status_by_pk\""
	InsertSaasAccount               *SaasAccountMutationResponse            "json:\"insert_saas_account\" graphql:\"insert_saas_account\""
	InsertSaasAccountOne            *SaasAccount                            "json:\"insert_saas_account_one\" graphql:\"insert_saas_account_one\""
	InsertSaasAddress               *SaasAddressMutationResponse            "json:\"insert_saas_address\" graphql:\"insert_saas_address\""
	InsertSaasAddressOne            *SaasAddress                            "json:\"insert_saas_address_one\" graphql:\"insert_saas_address_one\""
	InsertSaasMembership            *SaasMembershipMutationResponse         "json:\"insert_saas_membership\" graphql:\"insert_saas_membership\""
	InsertSaasMembershipOne         *SaasMembership                         "json:\"insert_saas_membership_one\" graphql:\"insert_saas_membership_one\""
	InsertSaasRole                  *SaasRoleMutationResponse               "json:\"insert_saas_role\" graphql:\"insert_saas_role\""
	InsertSaasRoleOne               *SaasRole                               "json:\"insert_saas_role_one\" graphql:\"insert_saas_role_one\""
	InsertSubscriptionActivePlan    *SubscriptionActivePlanMutationResponse "json:\"insert_subscription_active_plan\" graphql:\"insert_subscription_active_plan\""
	InsertSubscriptionActivePlanOne *SubscriptionActivePlan                 "json:\"insert_subscription_active_plan_one\" graphql:\"insert_subscription_active_plan_one\""
	InsertSubscriptionCustomer      *SubscriptionCustomerMutationResponse   "json:\"insert_subscription_customer\" graphql:\"insert_subscription_customer\""
	InsertSubscriptionCustomerOne   *SubscriptionCustomer                   "json:\"insert_subscription_customer_one\" graphql:\"insert_subscription_customer_one\""
	InsertSubscriptionEvent         *SubscriptionEventMutationResponse      "json:\"insert_subscription_event\" graphql:\"insert_subscription_event\""
	InsertSubscriptionEventOne      *SubscriptionEvent                      "json:\"insert_subscription_event_one\" graphql:\"insert_subscription_event_one\""
	InsertSubscriptionPlan          *SubscriptionPlanMutationResponse       "json:\"insert_subscription_plan\" graphql:\"insert_subscription_plan\""
	InsertSubscriptionPlanOne       *SubscriptionPlan                       "json:\"insert_subscription_plan_one\" graphql:\"insert_subscription_plan_one\""
	InsertSubscriptionStatus        *SubscriptionStatusMutationResponse     "json:\"insert_subscription_status\" graphql:\"insert_subscription_status\""
	InsertSubscriptionStatusOne     *SubscriptionStatus                     "json:\"insert_subscription_status_one\" graphql:\"insert_subscription_status_one\""
	SubscriptionCancel              *CancelSubscriptionOutput               "json:\"subscription_cancel\" graphql:\"subscription_cancel\""
	SubscriptionChange              *ChangeSubscriptionOutput               "json:\"subscription_change\" graphql:\"subscription_change\""
	SubscriptionCreate              *CreateSubscriptionOutput               "json:\"subscription_create\" graphql:\"subscription_create\""
	SubscriptionInit                *InitSubscriptionOutput                 "json:\"subscription_init\" graphql:\"subscription_init\""
	SubscriptionRetry               *RetrySubscriptionOutput                "json:\"subscription_retry\" graphql:\"subscription_retry\""
	UpdateSaasAccount               *SaasAccountMutationResponse            "json:\"update_saas_account\" graphql:\"update_saas_account\""
	UpdateSaasAccountByPk           *SaasAccount                            "json:\"update_saas_account_by_pk\" graphql:\"update_saas_account_by_pk\""
	UpdateSaasAddress               *SaasAddressMutationResponse            "json:\"update_saas_address\" graphql:\"update_saas_address\""
	UpdateSaasAddressByPk           *SaasAddress                            "json:\"update_saas_address_by_pk\" graphql:\"update_saas_address_by_pk\""
	UpdateSaasMembership            *SaasMembershipMutationResponse         "json:\"update_saas_membership\" graphql:\"update_saas_membership\""
	UpdateSaasMembershipByPk        *SaasMembership                         "json:\"update_saas_membership_by_pk\" graphql:\"update_saas_membership_by_pk\""
	UpdateSaasRole                  *SaasRoleMutationResponse               "json:\"update_saas_role\" graphql:\"update_saas_role\""
	UpdateSaasRoleByPk              *SaasRole                               "json:\"update_saas_role_by_pk\" graphql:\"update_saas_role_by_pk\""
	UpdateSubscriptionActivePlan    *SubscriptionActivePlanMutationResponse "json:\"update_subscription_active_plan\" graphql:\"update_subscription_active_plan\""
	UpdateSubscriptionCustomer      *SubscriptionCustomerMutationResponse   "json:\"update_subscription_customer\" graphql:\"update_subscription_customer\""
	UpdateSubscriptionCustomerByPk  *SubscriptionCustomer                   "json:\"update_subscription_customer_by_pk\" graphql:\"update_subscription_customer_by_pk\""
	UpdateSubscriptionEvent         *SubscriptionEventMutationResponse      "json:\"update_subscription_event\" graphql:\"update_subscription_event\""
	UpdateSubscriptionEventByPk     *SubscriptionEvent                      "json:\"update_subscription_event_by_pk\" graphql:\"update_subscription_event_by_pk\""
	UpdateSubscriptionPlan          *SubscriptionPlanMutationResponse       "json:\"update_subscription_plan\" graphql:\"update_subscription_plan\""
	UpdateSubscriptionPlanByPk      *SubscriptionPlan                       "json:\"update_subscription_plan_by_pk\" graphql:\"update_subscription_plan_by_pk\""
	UpdateSubscriptionStatus        *SubscriptionStatusMutationResponse     "json:\"update_subscription_status\" graphql:\"update_subscription_status\""
	UpdateSubscriptionStatusByPk    *SubscriptionStatus                     "json:\"update_subscription_status_by_pk\" graphql:\"update_subscription_status_by_pk\""
}
type MutationCreateSubscriptionCustomer struct {
	InsertSaasAccount *struct {
		AffectedRows int64 "json:\"affected_rows\" graphql:\"affected_rows\""
		Returning    []*struct {
			ID string "json:\"id\" graphql:\"id\""
		} "json:\"returning\" graphql:\"returning\""
	} "json:\"insert_saas_account\" graphql:\"insert_saas_account\""
}
type MutationSetSubscriptioStatus struct {
	UpdateSubscriptionStatus *struct {
		AffectedRows int64 "json:\"affected_rows\" graphql:\"affected_rows\""
		Returning    []*struct {
			IDAccount string "json:\"id_account\" graphql:\"id_account\""
			IsActive  bool   "json:\"is_active\" graphql:\"is_active\""
			Status    string "json:\"status\" graphql:\"status\""
		} "json:\"returning\" graphql:\"returning\""
	} "json:\"update_subscription_status\" graphql:\"update_subscription_status\""
}
type MutationAddSubscriptionEvent struct {
	InsertSubscriptionEvent *struct {
		AffectedRows int64 "json:\"affected_rows\" graphql:\"affected_rows\""
		Returning    []*struct {
			ID string "json:\"id\" graphql:\"id\""
		} "json:\"returning\" graphql:\"returning\""
	} "json:\"insert_subscription_event\" graphql:\"insert_subscription_event\""
}
type QueryGetAccountInfoForCreatingSubscription struct {
	SaasAccount []*struct {
		ID                   string "json:\"id\" graphql:\"id\""
		SubscriptionCustomer struct {
			StripeCustomer string "json:\"stripe_customer\" graphql:\"stripe_customer\""
		} "json:\"subscription_customer\" graphql:\"subscription_customer\""
		SubscriptionStatus struct {
			Status           string "json:\"status\" graphql:\"status\""
			SubscriptionPlan struct {
				StripeCode *string "json:\"stripe_code\" graphql:\"stripe_code\""
			} "json:\"subscription_plan\" graphql:\"subscription_plan\""
		} "json:\"subscription_status\" graphql:\"subscription_status\""
	} "json:\"saas_account\" graphql:\"saas_account\""
}
type QueryGetRoleForUserAndAccount struct {
	SaasMembership []*struct {
		IDRole string "json:\"id_role\" graphql:\"id_role\""
	} "json:\"saas_membership\" graphql:\"saas_membership\""
}
type QueryGetStripeSubscription struct {
	SubscriptionStatus []*struct {
		StripeSubscriptionID *string "json:\"stripe_subscription_id\" graphql:\"stripe_subscription_id\""
	} "json:\"subscription_status\" graphql:\"subscription_status\""
}
type QueryGetAccountFromSubscription struct {
	SubscriptionStatus []*struct {
		IDAccount string "json:\"id_account\" graphql:\"id_account\""
	} "json:\"subscription_status\" graphql:\"subscription_status\""
}

const CreateSubscriptionCustomerDocument = `mutation CreateSubscriptionCustomer ($name: String!, $id_plan: String!, $id_user: String!, $stripe_customer: String!, $status: String!, $id_role: String!) {
	insert_saas_account(objects: {name:$name,subscription_status:{data:{status:$status,id_plan:$id_plan}},subscription_customer:{data:{stripe_customer:$stripe_customer}},saas_memberships:{data:{id_role:$id_role,id_user:$id_user}}}) {
		affected_rows
		returning {
			id
		}
	}
}
`

func (c *Client) CreateSubscriptionCustomer(ctx context.Context, name string, idPlan string, idUser string, stripeCustomer string, status string, idRole string, interceptors ...clientv2.RequestInterceptor) (*MutationCreateSubscriptionCustomer, error) {
	vars := map[string]interface{}{
		"name":            name,
		"id_plan":         idPlan,
		"id_user":         idUser,
		"stripe_customer": stripeCustomer,
		"status":          status,
		"id_role":         idRole,
	}

	var res MutationCreateSubscriptionCustomer
	if err := c.Client.Post(ctx, "CreateSubscriptionCustomer", CreateSubscriptionCustomerDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const SetSubscriptioStatusDocument = `mutation SetSubscriptioStatus ($status: String!, $is_active: Boolean!, $accountId: uuid!, $stripe_subscription_id: String!) {
	update_subscription_status(where: {id_account:{_eq:$accountId}}, _set: {status:$status,is_active:$is_active,stripe_subscription_id:$stripe_subscription_id}) {
		affected_rows
		returning {
			id_account
			is_active
			status
		}
	}
}
`

func (c *Client) SetSubscriptioStatus(ctx context.Context, status string, isActive bool, accountID string, stripeSubscriptionID string, interceptors ...clientv2.RequestInterceptor) (*MutationSetSubscriptioStatus, error) {
	vars := map[string]interface{}{
		"status":                 status,
		"is_active":              isActive,
		"accountId":              accountID,
		"stripe_subscription_id": stripeSubscriptionID,
	}

	var res MutationSetSubscriptioStatus
	if err := c.Client.Post(ctx, "SetSubscriptioStatus", SetSubscriptioStatusDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const AddSubscriptionEventDocument = `mutation AddSubscriptionEvent ($type: String!, $data: jsonb!) {
	insert_subscription_event(objects: {data:$data,type:$type}) {
		affected_rows
		returning {
			id
		}
	}
}
`

func (c *Client) AddSubscriptionEvent(ctx context.Context, typeArg string, data map[string]interface{}, interceptors ...clientv2.RequestInterceptor) (*MutationAddSubscriptionEvent, error) {
	vars := map[string]interface{}{
		"type": typeArg,
		"data": data,
	}

	var res MutationAddSubscriptionEvent
	if err := c.Client.Post(ctx, "AddSubscriptionEvent", AddSubscriptionEventDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetAccountInfoForCreatingSubscriptionDocument = `query GetAccountInfoForCreatingSubscription ($id: uuid!) {
	saas_account(where: {id:{_eq:$id}}) {
		id
		subscription_customer {
			stripe_customer
		}
		subscription_status {
			status
			subscription_plan {
				stripe_code
			}
		}
	}
}
`

func (c *Client) GetAccountInfoForCreatingSubscription(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountInfoForCreatingSubscription, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var res QueryGetAccountInfoForCreatingSubscription
	if err := c.Client.Post(ctx, "GetAccountInfoForCreatingSubscription", GetAccountInfoForCreatingSubscriptionDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetRoleForUserAndAccountDocument = `query GetRoleForUserAndAccount ($user: String!, $account: uuid!) {
	saas_membership(where: {id_user:{_eq:$user},id_account:{_eq:$account}}) {
		id_role
	}
}
`

func (c *Client) GetRoleForUserAndAccount(ctx context.Context, user string, account string, interceptors ...clientv2.RequestInterceptor) (*QueryGetRoleForUserAndAccount, error) {
	vars := map[string]interface{}{
		"user":    user,
		"account": account,
	}

	var res QueryGetRoleForUserAndAccount
	if err := c.Client.Post(ctx, "GetRoleForUserAndAccount", GetRoleForUserAndAccountDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetStripeSubscriptionDocument = `query GetStripeSubscription ($id_account: uuid!) {
	subscription_status(where: {id_account:{_eq:$id_account}}) {
		stripe_subscription_id
	}
}
`

func (c *Client) GetStripeSubscription(ctx context.Context, idAccount string, interceptors ...clientv2.RequestInterceptor) (*QueryGetStripeSubscription, error) {
	vars := map[string]interface{}{
		"id_account": idAccount,
	}

	var res QueryGetStripeSubscription
	if err := c.Client.Post(ctx, "GetStripeSubscription", GetStripeSubscriptionDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetAccountFromSubscriptionDocument = `query GetAccountFromSubscription ($stripe_subscription_id: String!) {
	subscription_status(where: {stripe_subscription_id:{_eq:$stripe_subscription_id}}) {
		id_account
	}
}
`

func (c *Client) GetAccountFromSubscription(ctx context.Context, stripeSubscriptionID string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountFromSubscription, error) {
	vars := map[string]interface{}{
		"stripe_subscription_id": stripeSubscriptionID,
	}

	var res QueryGetAccountFromSubscription
	if err := c.Client.Post(ctx, "GetAccountFromSubscription", GetAccountFromSubscriptionDocument, &res, vars, interceptors...); err != nil {
		return nil, err
	}

	return &res, nil
}
