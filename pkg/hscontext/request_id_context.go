package hscontext

import "context"

func WithRequestIDValue(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, REQUEST_ID, requestId)
}

func RequestIDValue(ctx context.Context) string {
	return ctx.Value(REQUEST_ID).(string)
}
