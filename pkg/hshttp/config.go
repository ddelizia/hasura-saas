package hshttp

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigAccountIdHeader() string {
	return env.GetString("hshttp.headerNames.accountId")
}

func ConfigJWTHeader() string {
	return env.GetString("hshttp.headerNames.jwt")
}
