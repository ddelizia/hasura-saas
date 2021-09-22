package hshttp

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func SetHaderOnRequest(req *http.Request, header string, value string) {
	req.Header.Set(header, value)
	logrus.WithField(header, value).Debug("setting request header")
}

func SetHaderOnResponse(res http.ResponseWriter, header string, value string) {
	res.Header().Set(header, value)
	logrus.WithField(header, value).Debug("setting response header")
}

func SetCors(res http.ResponseWriter) {
	SetHaderOnResponse(res, "Access-Control-Allow-Origin", "http:localhost:8000")
	SetHaderOnResponse(res, "Access-Control-Allow-Methods", "OPTIONS,HEAD,GET,POST,PUT,DELETE")
	SetHaderOnResponse(res, "Access-Control-Allow-Headers", "X-Requested-With,X-Auth-Token,Content-Type,Content-Length,Authorization,Access-Control-Allow-Headers,Accept")
}

// SetSslRedirect Update the headers to allow for SSL redirection
func SetSslRedirect(req *http.Request, u *url.URL) {
	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme
	SetHaderOnRequest(req, "X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = u.Host
}

func AccountHeaderName() string {
	return ConfigAccountIdHeader()
}

func JwtHeaderName() string {
	return ConfigJWTHeader()
}

func CleanHasuraSaasHeaders(r *http.Request) {
	for k := range r.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-hasura-saas") {
			r.Header.Del(k)
		}
	}
}
