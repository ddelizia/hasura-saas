package rp

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type ReverseProxyMock struct {
	mock.Mock
}

func (m *ReverseProxyMock) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.Called(res, req)
}

func NewReverseProxyMock() http.Handler {
	return new(ReverseProxyMock)
}
