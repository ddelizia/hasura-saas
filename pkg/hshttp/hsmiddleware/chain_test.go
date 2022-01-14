package hsmiddleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/ddelizia/hasura-saas/pkg/hshttp/hsmiddleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

var _ = Describe("Chain", func() {

	logrus.SetOutput(ioutil.Discard)

	It("should chain multiple middlewares", func() {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mockHandlerParent := func(w http.ResponseWriter, r *http.Request) {}

		middleware1Visited := false
		middlaware1 := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				middleware1Visited = true
				next.ServeHTTP(w, r)
			}
		}

		middleware2Visited := false
		middlaware2 := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				middleware2Visited = true
				next.ServeHTTP(w, r)
			}
		}

		// When
		hsmiddleware.Chain(mockHandlerParent, middlaware1, middlaware2)(w, req)

		// Then
		Expect(middleware1Visited).To(BeTrue(), "middleware 1 has not been bisited")
		Expect(middleware2Visited).To(BeTrue(), "middleware 2 has not been bisited")
	})

	It("should load in order the middlewares", func() {
		// Given
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mockHandlerParent := func(w http.ResponseWriter, r *http.Request) {}

		visited := []string{}
		middlaware1 := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				visited = append(visited, "0")
				next.ServeHTTP(w, r)
			}
		}

		middlaware2 := func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				visited = append(visited, "1")
				next.ServeHTTP(w, r)
			}
		}

		// When
		hsmiddleware.Chain(mockHandlerParent, middlaware1, middlaware2)(w, req)

		// Then
		Expect(visited).To(Equal([]string{"0", "1"}), "middleware has not the same order")
	})

})
