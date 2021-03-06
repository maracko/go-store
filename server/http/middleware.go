package http

import (
	"log"
	"net/http"

	"github.com/maracko/go-store/errors"
	"github.com/maracko/go-store/server/http/helpers"
)

type middleware func(http.HandlerFunc) http.HandlerFunc

var commonMiddleware = []middleware{
	logMiddleWare,
	jsonHeader,
	authMiddleWare,
}

func logMiddleWare(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL, r.Host)
		h.ServeHTTP(w, r)
	})

}

func authMiddleWare(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if key != "" && header != key {
			helpers.JSONEncode(w, errors.Unauthorized("invalid key"))
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func jsonHeader(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func multipleMiddleware(h http.HandlerFunc, m ...middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}
