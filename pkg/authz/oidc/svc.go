package oidc

import "github.com/ddelizia/hasura-saas/pkg/authz"

type service struct {
	RoleGetterOidc
	AccountGetterOidc
	UserGetterOidc
	authz.AuthInfoGetter
}

func NewService() authz.Service {
	roleGetter := RoleGetterOidc{}
	accountGetter := AccountGetterOidc{}
	userGetter := UserGetterOidc{}
	return &service{
		RoleGetterOidc:    roleGetter,
		AccountGetterOidc: accountGetter,
		UserGetterOidc:    userGetter,
		AuthInfoGetter: &authz.AuthInfoGetterImpl{
			UserGetter:    &userGetter,
			AccountGetter: &accountGetter,
			RoleGetter:    &roleGetter,
		},
	}
}
