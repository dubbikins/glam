package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/dubbikins/glam/logging"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (recorder *StatusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return recorder.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	// r.ResponseWriter.WriteHeader(status)
}

func WithLogging(logger logging.Loggable) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &StatusRecorder{
				ResponseWriter: w,
				Status:         200,
			}
			next.ServeHTTP(recorder, r)
			logger.Info(fmt.Sprintf("%s %s %s => %s", logging.Cyan("HTTP"), logging.Purple(r.Method), logging.Gray(r.RequestURI), PickColor(fmt.Sprintf("%d", recorder.Status), recorder.Status)))
		})
	}
}

func PickColor(text string, status int) string {
	switch {
	case status < 300:
		return logging.Green(text)
	case status < 400:
		return logging.Cyan(text)
	default:
		return logging.Red(text)
	}
}
