package hsmiddleware

import (
	"net/http"
	"time"

	"github.com/ddelizia/hasura-saas/pkg/hscontext"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func LogRequest() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			requestID := uuid.NewString()
			newCtx := hscontext.WithRequestIDValue(r.Context(), requestID)
			r = r.WithContext(newCtx)
			logrus.WithContext(r.Context()).WithFields(logrus.Fields{
				"remote": r.RemoteAddr,
				"method": r.Method,
				"url":    r.URL,
			}).Debug("processing request")
			next.ServeHTTP(w, r)
			logrus.WithContext(r.Context()).WithFields(logrus.Fields{
				"responseTimeNS": time.Since(startTime).Nanoseconds(),
			}).Debug("response time")
		})
	}
}
