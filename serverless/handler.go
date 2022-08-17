package serverless

import (
	"context"
	"net/http"

	"github.com/dubbikins/glam"
)

type Requestable interface {
	Request() *http.Request
}

type Serverless[Request Requestable, Response any] struct {
	Router *glam.Router
	Domain string
	Writer http.ResponseWriter
}

func New[Request Requestable, Response any](router *glam.Router, writer http.ResponseWriter) *Serverless[Request, Response] {
	return &Serverless[Request, Response]{
		Router: router,
		Writer: writer,
	}
}

// events.ALBTargetGroupResponse
func (s *Serverless[Request, Response]) Handler() func(ctx context.Context, request Request) (Response, error) {
	return func(ctx context.Context, request Request) (Response, error) {
		s.Router.ServeHTTP(s.Writer, request.Request())
		return s.Writer.(Response), nil
	}
}

/**
Example: AWS Lambda Handler with ALB Target
func handler() {
	router := glam.NewRouter()
	//...setup router
	w := aws.NewResponseWriter(false) //true if Multivalue Header
	server := serverless.New[*aws.ALBTargetGroupRequest, *aws.ALBResponseWriter](router, w)
	lambda.Start(server.Handler())
}
**/
