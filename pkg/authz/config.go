package authz

import "github.com/ddelizia/hasura-saas/pkg/env"

func ConfigAnonymousUser() string {
	return env.GetString("authz.users.anonymous")
}

func ConfigAnonymousRole() string {
	return env.GetString("authz.roles.anonymous")
}

func ConfigAccountOwnerRole() string {
	return env.GetString("authz.roles.accountOwner")
}

func ConfigLoggedInRole() string {
	return env.GetString("authz.roles.loggedIn")
}

func ConfigAnonymousAccount() string {
	return env.GetString("authz.accounts.anonymous")
}
