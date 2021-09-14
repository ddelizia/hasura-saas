package gqlreq

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigGraphqlURL() string {
	return env.GetString("graphql.url")
}

func ConfigAdminSecret() string {
	return env.GetString("graphql.hasura.adminSecret")
}

func ConfigHasuraUserIdHeader() string {
	return env.GetString("graphql.hasura.headerNames.userId")
}

func ConfigHasuraRoleHeader() string {
	return env.GetString("graphql.hasura.headerNames.role")
}

func ConfigHasuraAccountHeader() string {
	return env.GetString("graphql.hasura.headerNames.accountId")
}
