package hscontext

type ContextKey int

const (
	ACTION_SESSION_VARIABLES ContextKey = 100001
	ACTION_DATA              ContextKey = 100002
	REQUEST_ID               ContextKey = 100003
	AUTHZ_INFO               ContextKey = 100004
)
