package email

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigFrom() string {
	return env.GetString("email.from")
}
