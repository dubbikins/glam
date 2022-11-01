package glam

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dubbikins/glam/logging"
)

type UrlParameter struct {
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})
var logger = logging.Logger

type Router struct {
	root            *Node
	NotFoundHandler http.Handler
}

func (r *Router) Root() *Node {
	return r.root
}

// NewRouter creates light weight composible router that implements the http.Handler interface
// You can add handlers for specific paths, match on key words or regular expressions.
func NewRouter() *Router {
	r := &Router{
		NotFoundHandler: defaultNotFoundHandler,
	}
	r.root = NewRoot(r)
	return r
}

// Handler adds a handler function to the router for a path and method combination
func (c *Router) Handle(path, method string, handler http.HandlerFunc) {
	path = c.root.Name + path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	paths := strings.Split(path, "/")
	for len(paths) > 0 && paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}
	c.root.insertHandler(paths, method, handler)
}

// Get add a GET method handler for the speficied path
func (c *Router) Get(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodGet, handler)
}

// Patch add a PATCH method handler for the speficied path
func (c *Router) Patch(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPatch, handler)
}

// Put add a PUT method handler for the speficied path
func (c *Router) Put(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPut, handler)
}

// Post add a POST method handler for the speficied path
func (c *Router) Post(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodPost, handler)
}

// Delete add a DELETE method handler for the speficied path
func (c *Router) Delete(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodDelete, handler)
}

// Head adds a HEAD method handler for the speficied path
func (c *Router) Head(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodHead, handler)
}

// Options adds an OPTIONS method handler for the speficied path
func (c *Router) Options(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodOptions, handler)
}

// Connect adds a CONNECT method handler for the speficied path
func (c *Router) Connect(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodConnect, handler)
}

// Trace adds a TRACE method handler for the speficied path
func (c *Router) Trace(path string, handler http.HandlerFunc) {
	c.Handle(path, http.MethodTrace, handler)
}

// Use adds 1 or more middleware to the router.
// The middleware function will be applied to all handlers associated with this router
func (r *Router) Use(middleware ...Middleware) {
	r.root.insertMiddleware([]string{}, middleware)
}

// UseAt adds 1 or more middleware to a path in the router.
// The middleware function will be applied to all handlers with a prefix matching this path
// func (r *Router) UseAt(path string, middleware ...Middleware) {
// 	path = r.root.Name + path
// 	if strings.HasPrefix(path, "/") {
// 		path = path[1:]
// 	}
// 	paths := strings.Split(path, "/")
// 	r.root.insertMiddleware(paths, middleware)
// }

// Router adds a set of routes to the router at a specified path.
// This is convinient for paths the share a common prefix reduing the need to specify the entire path for every handler.
func (router *Router) Routes(path string, setRoutes func(r *Router)) {
	newRouter := NewRouter()
	setRoutes(newRouter)
	router.Mount(path, newRouter)
}

// NotFound replaces the routers default NotFound handler function.
func (r *Router) NotFound(handler http.HandlerFunc) {
	r.NotFoundHandler = handler
}

// NotFoundAt adds a not found handler to a path
// This allows for multiple NotFound Handlers in the routing tree
// func (r *Router) NotFoundAt(path string, handler http.HandlerFunc, applyMiddleware bool) {
// 	method := "NOTFOUND"
// 	if applyMiddleware {
// 		method = "NOTFOUNDAPPLYMIDDLEWARE"
// 	}
// 	r.Handle(path, method, handler)
// }

// Mount converts a top-level router into a top-level node in the parent router with a specified prefix
func (parent *Router) Mount(prefix string, r *Router) {
	parent.root.InsertNodeAt(prefix, r.root)
}

// ServeHTTP is the Router's implementation of the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("%s %s", logging.Green(r.Method), logging.Purple(r.URL.Path)))
	router.root.ServeHTTP(w, r)
}
