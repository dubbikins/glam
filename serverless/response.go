package serverless
import (
	"net/http"
	"strings"
)
type ALBStringResponseWriter struct {
	StatusCode int `json:"statusCode"`
	StatusDescription string `json:"statusDescription"`
	IsBase64Encoded bool `json:"isBase64Encoded"`
	MultiValueHeaders http.Header  `json:"multiValueHeaders,omitempty"`
	Headers map[string] string `json:"headers,omitempty"`
	Body string `json:"body"`
	IsMultiValue bool `json:"-"`
}

func (w *ALBStringResponseWriter) Write(result []byte) (int, error) {
	if w.StatusCode == http.StatusNotImplemented {
		w.WriteHeader(http.StatusOK)
	}
	w.Body = string(result)
	return len(result), nil
}

func (w *ALBStringResponseWriter) WriteHeader(httpStatus int) {
	w.StatusCode = httpStatus
	w.StatusDescription = http.StatusText(httpStatus)
	if !w.IsMultiValue {
		for header, value := range w.MultiValueHeaders {
			w.Headers[header] = strings.Join(value, ";")
		}
		w.MultiValueHeaders = nil
	}
}

func (w *ALBStringResponseWriter) Header() http.Header {
	return w.MultiValueHeaders
}

func NewResponseWriter( isMultiValue bool) http.ResponseWriter {
	return &ALBStringResponseWriter{
		StatusCode: http.StatusNotImplemented,
		StatusDescription: "Not Implemented",
		Body: "",
		MultiValueHeaders: NewHttpHeaders(),
		Headers: make(map[string]string),
		IsBase64Encoded: false,
		IsMultiValue: isMultiValue,
	}
}

func NewHttpHeaders() http.Header {
	return make(http.Header)
}