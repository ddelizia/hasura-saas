package subscription

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigWebhookSecret() string {
	return env.GetString("subscription.stripe.webhookSecret")
}

func ConfigDomain() string {
	return env.GetString("subscription.stripe.domain")
}

func ConfigWebhookListenAddress() string {
	return env.GetString("subscription.server.listenAddress")
}

func ConfigApiKey() string {
	return env.GetString("subscription.stripe.apiKey")
}
