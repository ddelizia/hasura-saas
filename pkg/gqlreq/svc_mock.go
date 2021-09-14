package gqlreq

type ServiceMock struct {
	ExecuterMock
	HeaderInfoGetterMock
	SessionInfoGetterMock
}

func NewServiceMock() Service {
	return new(ServiceMock)
}
