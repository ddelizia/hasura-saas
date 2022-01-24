package hsmiddleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}

	// otherwise nest the handlerfuncs
	return m[0](Chain(f, m[1:cap(m)]...))
}
