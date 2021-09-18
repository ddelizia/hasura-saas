package hshttp

import (
	"net/http"

	"github.com/joomcode/errorx"
)

var Mapping map[string]int = map[string]int{
	errorx.IllegalArgument.FullName():      http.StatusBadRequest,
	errorx.IllegalState.FullName():         http.StatusInternalServerError,
	errorx.IllegalFormat.FullName():        http.StatusBadRequest,
	errorx.InitializationFailed.FullName(): http.StatusInternalServerError,
	errorx.DataUnavailable.FullName():      http.StatusNotFound,
	errorx.UnsupportedOperation.FullName(): http.StatusBadRequest,
	errorx.RejectedOperation.FullName():    http.StatusBadRequest,
	errorx.Interrupted.FullName():          http.StatusBadRequest,
	errorx.AssertionFailed.FullName():      http.StatusBadRequest,
	errorx.InternalError.FullName():        http.StatusInternalServerError,
	errorx.ExternalError.FullName():        http.StatusInternalServerError,
	errorx.ConcurrentUpdate.FullName():     http.StatusInternalServerError,
	errorx.TimeoutElapsed.FullName():       http.StatusRequestTimeout,
	errorx.NotImplemented.FullName():       http.StatusInternalServerError,
	errorx.UnsupportedVersion.FullName():   http.StatusInternalServerError,
}

func errorToStatusMapping(e error) int {
	typeName := errorx.GetTypeName(e)

	return (Mapping[typeName])
}

func WriteError(res http.ResponseWriter, e error) {
	status := errorToStatusMapping(e)
	if status != 0 {
		WriteResponse(res, status, struct {
			Error string `json:"error"`
		}{Error: e.Error()})
	} else {
		WriteResponse(res, 500, struct{ Error string }{Error: e.Error()})
	}
}
