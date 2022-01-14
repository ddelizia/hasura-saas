package hsmiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func RequestTimeout(timeout int) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logrus.WithContext(r.Context()).Debug("adding")
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout))
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
