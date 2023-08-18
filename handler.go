package glam

import (
	"net/http"
	"regexp"
)

type handler struct {
	methodHandlers map[string]http.Handler
	regexChildren  map[string]*regexHandler
	staticChildren map[string]*staticHandler
	strictChildren map[string]*strictHandler
	paramChild     *paramHandler
	middleware     []Middleware //24
	notFound       http.Handler
}

func (h *handler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

type regexHandler struct {
	handler
	matcher *regexp.Regexp
}

func (h *regexHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

type strictHandler struct {
	handler
}

func (h *strictHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

type paramHandler struct {
	handler
}

func (h *paramHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

type staticHandler struct {
	handler
}

func (h *staticHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}
