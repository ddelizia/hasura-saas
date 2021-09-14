package hshttp

import (
	"encoding/json"
	"net/http"
)

func WriteResponse(res http.ResponseWriter, statusCode int, data interface{}) {
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
}
