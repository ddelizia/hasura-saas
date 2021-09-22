package authz

type ServiceMock struct {
	RoleGetterMock
	UserGetterMock
	AccountGetterMock
	AuthInfoGetterMock
}

func NewServiceMock() Service {
	return new(ServiceMock)
}
