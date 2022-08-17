package middleware

import (
	"net/http"
	"net/url"
)

func RestoreURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL, _ = url.Parse(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
