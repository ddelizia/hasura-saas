package email

import "github.com/stretchr/testify/mock"

type ServiceMock struct {
	mock.Mock
}

func NewServiceMock() Service {
	return new(ServiceMock)
}

func (m *ServiceMock) SendEmail(emails []string, subject, template string, data interface{}) error {
	args := m.Called(emails, subject, template, data)
	return args.Error(0)
}
