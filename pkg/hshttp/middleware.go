package hshttp

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func MiddlewareChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](MiddlewareChain(f, m[1:cap(m)]...))
}

func MiddlewareSetContentTypeApplicationJson(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetHaderOnResponse(w, "Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func MiddlewareLogRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		requestID := uuid.NewString()
		logrus.WithFields(logrus.Fields{
			"remote":    r.RemoteAddr,
			"method":    r.Method,
			"url":       r.URL,
			"requestID": requestID,
		}).Debug("processing request")
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"responseTimeNS": time.Since(startTime).Nanoseconds(),
			"requestID":      requestID,
		}).Debug("response time")
	})
}

func MiddlewareRequestContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("changing request context")
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
