package gqlsdk

import (
	"context"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() Service {
	return &ServiceMock{}
}

func (m *ServiceMock) CreateSubscriptionCustomer(ctx context.Context, name string, idPlan string, idUser string, stripeCustomer string, status string, idRole string, interceptors ...clientv2.RequestInterceptor) (*MutationCreateSubscriptionCustomer, error) {
	args := m.Called(ctx, name, idPlan, idUser, stripeCustomer, status, idRole, interceptors)
	return args.Get(0).(*MutationCreateSubscriptionCustomer), args.Error(1)
}

func (m *ServiceMock) GetAccountInfoForCreatingSubscription(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountInfoForCreatingSubscription, error) {
	args := m.Called(ctx, id, interceptors)
	return args.Get(0).(*QueryGetAccountInfoForCreatingSubscription), args.Error(1)
}

func (m *ServiceMock) SetSubscriptioStatus(ctx context.Context, status string, isActive bool, accountID string, stripeSubscriptionID string, idPlan string, interceptors ...clientv2.RequestInterceptor) (*MutationSetSubscriptioStatus, error) {
	args := m.Called(ctx, status, isActive, accountID, stripeSubscriptionID, idPlan, interceptors)
	return args.Get(0).(*MutationSetSubscriptioStatus), args.Error(1)
}

func (m *ServiceMock) GetStripeSubscription(ctx context.Context, idAccount string, interceptors ...clientv2.RequestInterceptor) (*QueryGetStripeSubscription, error) {
	args := m.Called(ctx, idAccount, interceptors)
	return args.Get(0).(*QueryGetStripeSubscription), args.Error(1)
}

func (m *ServiceMock) AddSubscriptionEvent(ctx context.Context, typeArg string, data map[string]interface{}, interceptors ...clientv2.RequestInterceptor) (*MutationAddSubscriptionEvent, error) {
	args := m.Called(ctx, typeArg, data, interceptors)
	return args.Get(0).(*MutationAddSubscriptionEvent), args.Error(1)
}

func (m *ServiceMock) GetRoleForUserAndAccount(ctx context.Context, user string, account string, interceptors ...clientv2.RequestInterceptor) (*QueryGetRoleForUserAndAccount, error) {
	args := m.Called(ctx, user, account, interceptors)
	return args.Get(0).(*QueryGetRoleForUserAndAccount), args.Error(1)
}

func (m *ServiceMock) GetAccountFromSubscription(ctx context.Context, stripeSubscriptionID string, interceptors ...clientv2.RequestInterceptor) (*QueryGetAccountFromSubscription, error) {
	args := m.Called(ctx, stripeSubscriptionID, interceptors)
	return args.Get(0).(*QueryGetAccountFromSubscription), args.Error(1)
}

func (m *ServiceMock) GetStripePlanFromPlan(ctx context.Context, id string, interceptors ...clientv2.RequestInterceptor) (*QueryGetStripePlanFromPlan, error) {
	args := m.Called(ctx, id, interceptors)
	return args.Get(0).(*QueryGetStripePlanFromPlan), args.Error(1)
}

func (m *ServiceMock) GetPlanFromStripePlan(ctx context.Context, stripeCode string, interceptors ...clientv2.RequestInterceptor) (*QueryGetPlanFromStripePlan, error) {
	args := m.Called(ctx, stripeCode, interceptors)
	return args.Get(0).(*QueryGetPlanFromStripePlan), args.Error(1)
}

func (m *ServiceMock) SetAccountForUser(ctx context.Context, idAccount string, idUser string, interceptors ...clientv2.RequestInterceptor) (*MutationSetAccountForUser, error) {
	args := m.Called(ctx, idAccount, idUser, interceptors)
	return args.Get(0).(*MutationSetAccountForUser), args.Error(1)
}

func (m *ServiceMock) GetCurrentAccount(ctx context.Context, userID string, interceptors ...clientv2.RequestInterceptor) (*QueryGetCurrentAccount, error) {
	args := m.Called(ctx, userID, interceptors)
	return args.Get(0).(*QueryGetCurrentAccount), args.Error(1)
}
