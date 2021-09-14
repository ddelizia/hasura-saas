package hstest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockClient struct {
	mock.Mock
}

var DefaultMockClient = &MockClient{}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func BuildResponseBodyForMock(json string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(json)))
}
