package authz

import "github.com/ddelizia/hasura-saas/pkg/gqlsdk"

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

func NewService(sdkService gqlsdk.Service) Service {
	roleGetterImpl := RoleGetterImpl{GraphQlSvc: sdkService}
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
