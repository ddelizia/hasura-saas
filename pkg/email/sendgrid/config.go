package sendgrid

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigSendgridApiKey() string {
	return env.GetString("email.sendgrid.apikey")
}