package hstest

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onsi/ginkgo"
)

func CrerateRequest(method string, path string, body io.Reader, vars map[string]string) *http.Request {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		ginkgo.Fail("Error executing the request")
	}

	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}

	return req
}
