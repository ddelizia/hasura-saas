package hscontext

import (
	"context"
)

// Create
func WithActionSessionValue(ctx context.Context, sessionVariables map[string]interface{}, data interface{}) context.Context {
	newCtx := context.WithValue(ctx, ACTION_SESSION_VARIABLES, sessionVariables)
	newCtx = context.WithValue(newCtx, ACTION_DATA, data)
	return newCtx
}

func ActionSessionVariablesValue(ctx context.Context) map[string]interface{} {
	return ctx.Value(ACTION_SESSION_VARIABLES).(map[string]interface{})
}

func ActionDataValue(ctx context.Context) interface{} {
	return ctx.Value(ACTION_DATA)
}
