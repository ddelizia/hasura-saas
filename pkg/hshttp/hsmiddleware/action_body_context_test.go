package hsmiddleware_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ddelizia/hasura-saas/pkg/hasura"
	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	"github.com/ddelizia/hasura-saas/pkg/hshttp/hsmiddleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

type InputData struct {
	InputKey string `json:"input_key"`
}

type InputPayloadPtr struct {
	hasura.BasePayload
	Input *InputData `json:"input"`
}

type InputPayload struct {
	hasura.BasePayload
	Input InputData `json:"input"`
}

var _ = Describe("ActionBodyContext", func() {

	logrus.SetOutput(ioutil.Discard)

	exampleBody := "{\"session_variables\": {\"varKey\": \"varValue\"}, \"input\": {\"input_key\":\"inputValue\"}}"

	It("should store data into context", func() {
		// Given
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(exampleBody))
		w := httptest.NewRecorder()
		nextHasBeenCalled := false
		var sessionVariables map[string]interface{}
		var inputData *InputData
		mockNext := func(w http.ResponseWriter, r *http.Request) {
			nextHasBeenCalled = true
			sessionVariables = hscontext.ActionSessionVariablesValue(r.Context())
			inputData = hscontext.ActionDataValue(r.Context()).(*InputData)
		}

		// When
		hsmiddleware.ActionBodyToContext(&InputPayloadPtr{})(mockNext)(w, req)

		// Then
		Expect(nextHasBeenCalled).To(BeTrue())
		Expect(sessionVariables).ToNot(BeNil())
		Expect(sessionVariables["varKey"]).To(Equal("varValue"))
		Expect(inputData).ToNot(BeNil())
		Expect(inputData.InputKey).To(Equal("inputValue"))
	})

	It("should store data into context and input is not a pointer", func() {
		// Given
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(exampleBody))
		w := httptest.NewRecorder()
		nextHasBeenCalled := false
		var sessionVariables map[string]interface{}
		var inputData InputData
		mockNext := func(w http.ResponseWriter, r *http.Request) {
			nextHasBeenCalled = true
			sessionVariables = hscontext.ActionSessionVariablesValue(r.Context())
			inputData = hscontext.ActionDataValue(r.Context()).(InputData)
		}

		// When
		hsmiddleware.ActionBodyToContext(InputPayload{})(mockNext)(w, req)

		// Then
		Expect(nextHasBeenCalled).To(BeTrue())
		Expect(sessionVariables).ToNot(BeNil())
		Expect(sessionVariables["varKey"]).To(Equal("varValue"))
		Expect(inputData).ToNot(BeNil())
		Expect(inputData.InputKey).To(Equal("inputValue"))
	})

	It("should be able to deal with multiple bodies", func() {
		// Given
		exampleNewBody := "{\"session_variables\": {\"varKey1\": \"varValue\"}, \"input\": {\"input_key\":\"inputValue2\"}}"
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(exampleBody))
		reqNew := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(exampleNewBody))
		w := httptest.NewRecorder()
		wNew := httptest.NewRecorder()
		nextHasBeenCalled := false
		var sessionVariables map[string]interface{}
		var inputData *InputData
		mockNext := func(w http.ResponseWriter, r *http.Request) {
			nextHasBeenCalled = true
			sessionVariables = hscontext.ActionSessionVariablesValue(r.Context())
			inputData = hscontext.ActionDataValue(r.Context()).(*InputData)
		}

		// When
		handler := hsmiddleware.ActionBodyToContext(&InputPayloadPtr{})(mockNext)
		handler(w, req)
		handler(wNew, reqNew)

		// Then
		Expect(nextHasBeenCalled).To(BeTrue())
		Expect(sessionVariables).ToNot(BeNil())
		Expect(sessionVariables["varKey1"]).To(Equal("varValue"))
		Expect(sessionVariables["varKey"]).To(BeNil())
		Expect(inputData).ToNot(BeNil())
		Expect(inputData.InputKey).To(Equal("inputValue2"))
	})

	It("should respond with errror when input is not valid json", func() {
		// Given
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("hjlkhjlkjhlkh"))
		w := httptest.NewRecorder()

		mockNext := func(w http.ResponseWriter, r *http.Request) {}

		// When
		hsmiddleware.ActionBodyToContext(&InputPayloadPtr{})(mockNext)(w, req)

		// Then
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		Expect(body).To(ContainSubstring("invalid payload for request"))
		Expect(resp.StatusCode).To(Equal(400), "it should be a bad request")
	})
})
