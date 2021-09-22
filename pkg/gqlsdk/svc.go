package gqlsdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
)

type Service interface {
	CreateSubscriptionCustomer(ctx context.Context, name string, idPlan string, idUser string, stripeCustomer string, status string, idRole string, interceptors ...clientv2.RequestInterceptor) (*MutationCreateSubscriptionCustomer, error)
	GetAccountInfoForCreatingSubscription(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountInfoForCreatingSubscription, error)
	SetSubscriptioStatus(ctx context.Context, status string, isActive bool, accountID string, stripeSubscriptionID string, idPlan string, interceptors ...clientv2.RequestInterceptor) (*MutationSetSubscriptioStatus, error)
	GetStripeSubscription(ctx context.Context, idAccount string, interceptors ...clientv2.RequestInterceptor) (*QueryGetStripeSubscription, error)
	AddSubscriptionEvent(ctx context.Context, typeArg string, data map[string]interface{}, interceptors ...clientv2.RequestInterceptor) (*MutationAddSubscriptionEvent, error)
	GetRoleForUserAndAccount(ctx context.Context, user string, account string, interceptors ...clientv2.RequestInterceptor) (*QueryGetRoleForUserAndAccount, error)
	GetAccountFromSubscription(ctx context.Context, stripeSubscriptionID string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountFromSubscription, error)
	GetStripePlanFromPlan(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*QueryGetStripePlanFromPlan, error)
	GetPlanFromStripePlan(ctx context.Context, stripeCode string, interceptors ...clientv2.RequestInterceptor) (*QueryGetPlanFromStripePlan, error)
	SetAccountForUser(ctx context.Context, idAccount string, idUser string, interceptors ...clientv2.RequestInterceptor) (*MutationSetAccountForUser, error)
	GetCurrentAccount(ctx context.Context, userID string, interceptors ...clientv2.RequestInterceptor) (*QueryGetCurrentAccount, error)
}
