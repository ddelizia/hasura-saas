package hsmiddleware

import (
	"net/http"
	"reflect"

	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	"github.com/ddelizia/hasura-saas/pkg/hshttp"
	"github.com/joomcode/errorx"
	"github.com/sirupsen/logrus"
)

// injects in the request context the session variables
func ActionBodyToContext(d interface{}) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logrus.Debug("parsing the body")
			var data interface{}
			dataType := reflect.TypeOf(d)
			if dataType.Kind() == reflect.Ptr {
				data = reflect.New(dataType.Elem()).Interface()
			} else {
				data = reflect.New(dataType).Interface()
			}
			err := hshttp.GetBody(r, data)
			if err != nil {
				hshttp.WriteError(w, errorx.IllegalArgument.Wrap(err, "invalid payload for request"))
				return
			}
			valueOfPayload := reflect.ValueOf(data).Elem()
			newCtx := hscontext.WithActionSessionValue(r.Context(), valueOfPayload.FieldByName("SessionVariables").Interface().(map[string]interface{}), valueOfPayload.FieldByName("Input").Interface())
			r = r.WithContext(newCtx)
			next.ServeHTTP(w, r)
		})
	}
}
