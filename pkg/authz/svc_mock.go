package authz

type ServiceMock struct {
	RoleGetterMock
	UserGetterMock
	AuthInfoGetterMock
}

func NewServiceMock() Service {
	return new(ServiceMock)
}
