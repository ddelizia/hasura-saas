package hscontext

import (
	"context"

	"github.com/ddelizia/hasura-saas/pkg/gqlreq"
)

func WithAuthzInfoValue(ctx context.Context, autzInfo *gqlreq.HeaderInfo) context.Context {
	return context.WithValue(ctx, AUTHZ_INFO, autzInfo)
}

func AuthzInfoValue(ctx context.Context) *gqlreq.HeaderInfo {
	return ctx.Value(AUTHZ_INFO).(*gqlreq.HeaderInfo)
}
