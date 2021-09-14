package authz

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigJwksUrl() string {
	return env.GetString("authz.jwks.url")
}

func ConfigAnonymousUser() string {
	return env.GetString("authz.users.anonymous")
}

func ConfigAnonymousRole() string {
	return env.GetString("authz.roles.anonymous")
}

func ConfigLoggedInRole() string {
	return env.GetString("authz.roles.loggedIn")
}

func ConfigAdminRole() string {
	return env.GetString("authz.roles.admin")
}

func ConfigAccountOwnerRole() string {
	return env.GetString("authz.roles.accountOwner")
}
