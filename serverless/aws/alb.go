package aws

import (
	"fmt"
	"net/http"
	"strings"
)

type ALBTargetGroupRequest struct {
	HTTPMethod                      string                       `json:"httpMethod"`
	Path                            string                       `json:"path"`
	QueryStringParameters           map[string]string            `json:"queryStringParameters,omitempty"`
	MultiValueQueryStringParameters map[string][]string          `json:"multiValueQueryStringParameters,omitempty"`
	Headers                         map[string]string            `json:"headers,omitempty"`
	MultiValueHeaders               map[string][]string          `json:"multiValueHeaders,omitempty"`
	RequestContext                  ALBTargetGroupRequestContext `json:"requestContext"`
	IsBase64Encoded                 bool                         `json:"isBase64Encoded"`
	Body                            string                       `json:"body,omitempty"`
}

func (event *ALBTargetGroupRequest) Request() *http.Request {
	path := event.Path
	if len(event.QueryStringParameters) > 0 {
		queryStringList := make([]string, len(event.QueryStringParameters))
		index := 0
		for key, value := range event.QueryStringParameters {
			queryStringList[index] = fmt.Sprintf("%s=%s", key, value)
			index++
		}
		queryString := strings.Join(queryStringList, "&")
		path += "?" + queryString
	}
	request, err := http.NewRequest(event.HTTPMethod, path, strings.NewReader(event.Body))
	if err != nil {
		panic("Error Creating Request from ALB Request Event")
	}
	return request
}

// ALBTargetGroupRequestContext contains the information to identify the load balancer invoking the lambda
type ALBTargetGroupRequestContext struct {
	ELB ELBContext `json:"elb"`
}

// ELBContext contains the information to identify the ARN invoking the lambda
type ELBContext struct {
	TargetGroupArn string `json:"targetGroupArn"` //nolint: stylecheck
}

type ALBResponseWriter struct {
	StatusCode        int               `json:"statusCode"`
	StatusDescription string            `json:"statusDescription"`
	IsBase64Encoded   bool              `json:"isBase64Encoded"`
	MultiValueHeaders http.Header       `json:"multiValueHeaders,omitempty"`
	Headers           map[string]string `json:"headers,omitempty"`
	Body              string            `json:"body"`
	IsMultiValue      bool              `json:"-"`
}

func (w *ALBResponseWriter) Write(result []byte) (int, error) {
	if w.StatusCode == http.StatusNotImplemented {
		w.WriteHeader(http.StatusOK)
	}
	w.Body = string(result)
	return len(result), nil
}

func (w *ALBResponseWriter) WriteHeader(httpStatus int) {
	w.StatusCode = httpStatus
	w.StatusDescription = http.StatusText(httpStatus)
	if !w.IsMultiValue {
		for header, value := range w.MultiValueHeaders {
			w.Headers[header] = strings.Join(value, ";")
		}
		w.MultiValueHeaders = nil
	}
}

func (w *ALBResponseWriter) Header() http.Header {
	return w.MultiValueHeaders
}

func NewResponseWriter(isMultiValue bool) http.ResponseWriter {
	return &ALBResponseWriter{
		StatusCode:        http.StatusNotImplemented,
		StatusDescription: "Not Implemented",
		Body:              "",
		MultiValueHeaders: NewHttpHeaders(),
		Headers:           make(map[string]string),
		IsBase64Encoded:   false,
		IsMultiValue:      isMultiValue,
	}
}

func NewHttpHeaders() http.Header {
	return make(http.Header)
}
