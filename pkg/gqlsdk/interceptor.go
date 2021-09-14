package gqlsdk

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/clientv2"
	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
)

func HasuraAdminInterceptor() clientv2.RequestInterceptor {
	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		hshttp.SetHaderOnRequest(req, gqlreq.HASURA_ADMIN_SECRET_HEADER_NAME, gqlreq.ConfigAdminSecret())

		return next(ctx, req, gqlInfo, res)
	}
}

func WithAdminRole() clientv2.RequestInterceptor {
	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraRoleHeader(), "admin")

		return next(ctx, req, gqlInfo, res)
	}
}

func AuthzHeadersInterceptor(accountId, userId, role string) clientv2.RequestInterceptor {
	return func(ctx context.Context, req *http.Request, gqlInfo *clientv2.GQLRequestInfo, res interface{}, next clientv2.RequestInterceptorFunc) error {
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraAccountHeader(), accountId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraUserIdHeader(), userId)
		hshttp.SetHaderOnRequest(req, gqlreq.ConfigHasuraRoleHeader(), role)

		return next(ctx, req, gqlInfo, res)
	}
}
