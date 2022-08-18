package glam

import (
	"net/http"
	"strings"
)

type UrlParameter struct {
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

type Router struct {
	root            *Node
	NotFoundHandler http.Handler
}

func (r *Router) Root() *Node {
	return r.root
}

func NewRouter() *Router {
	r := &Router{
		NotFoundHandler: defaultNotFoundHandler,
	}
	r.root = NewRoot(r)
	return r
}

func (c *Router) Handle(path, method string, handler http.HandlerFunc) {
	path = c.root.Name + path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	paths := strings.Split(path, "/")
	c.root.insertHandler(paths, method, handler)
}
func (c *Router) Get(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodGet, handler)
}

func (c *Router) Patch(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPatch, handler)
}

func (c *Router) Put(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPut, handler)
}

func (c *Router) Post(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPost, handler)
}

func (c *Router) Delete(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodDelete, handler)
}

func (c *Router) Head(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodHead, handler)
}

func (c *Router) Options(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodOptions, handler)
}

func (c *Router) Connect(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodConnect, handler)
}

func (c *Router) Trace(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodTrace, handler)
}
func (r *Router) Use(middleware ...Middleware) {

	r.root.insertMiddleware([]string{}, middleware)
}
func (r *Router) UseAt(path string, middleware ...Middleware) {
	path = r.root.Name + path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	paths := strings.Split(path, "/")
	r.root.insertMiddleware(paths, middleware)
}

func (router *Router) Router(path string, setRoutes func(r *Router)) {
	newRouter := NewRouter()
	setRoutes(newRouter)
	router.Mount(path, newRouter)
}

func (r *Router) NotFound(handler http.HandlerFunc) {
	r.NotFoundHandler = handler
}
func (r *Router) NotFoundAt(path string, handler http.HandlerFunc, applyMiddleware bool) {
	method := "NOTFOUND"
	if applyMiddleware {
		method = "NOTFOUNDAPPLYMIDDLEWARE"
	}
	r.Handle(path, method, handler)
}

func (parent *Router) Mount(prefix string, r *Router) {

	parent.root.InsertNodeAt(prefix, r.root)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.root.ServeHTTP(w, r)
}
