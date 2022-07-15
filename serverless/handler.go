package serverless

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/dubbikins/glam/router"
	"github.com/aws/aws-lambda-go/events"
)

type Serverless struct {
	Router *router.Router
	Domain string
}

func NewHandler(router *router.Router) *Serverless {
	return &Serverless{
		Router: router,
	}
}

func (s *Serverless) Response(event *events.ALBTargetGroupRequest, ctx context.Context) (events.ALBTargetGroupResponse, error) {
	w := NewResponseWriter(event.MultiValueHeaders != nil).(*ALBStringResponseWriter)
	r := s.RequestFromALBEvent(event)
	s.Router.ServeHTTP(w, r)

	return events.ALBTargetGroupResponse{
		Body: w.Body,
		StatusCode: w.StatusCode,
		StatusDescription: w.StatusDescription,
		IsBase64Encoded: w.IsBase64Encoded,
		Headers: w.Headers,
	}, nil
} 

func (s *Serverless) RequestFromALBEvent(event *events.ALBTargetGroupRequest) *http.Request {
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