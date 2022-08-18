package middleware

import (
	"net/http"
)

func WithContentType(content_type string) func(http.Handler) http.Handler {
	return WithHeader("Content-Type", content_type)
}
