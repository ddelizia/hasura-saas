package authz

import "github.com/ddelizia/hasura-saas/pkg/gqlreq"

type Service interface {
	RoleGetter
	UserGetter
	AuthInfoGetter
}

type service struct {
	RoleGetterImpl
	UserGetterImpl
	AuthInfoGetterImpl
}

func NewService() Service {
	roleGetterImpl := RoleGetterImpl{GraphQlSvc: gqlreq.NewService()}
	userGetterImpl := UserGetterImpl{}
	return &service{
		RoleGetterImpl: roleGetterImpl,
		UserGetterImpl: userGetterImpl,
		AuthInfoGetterImpl: AuthInfoGetterImpl{
			UserGetter: &userGetterImpl,
			RoleGetter: &roleGetterImpl,
		},
	}
}
