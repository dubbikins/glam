package glam

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dubbikins/glam/logging"
)

var logger = logging.Logger

type Router struct {
	root *node
}

func (r *Router) getRoot() *node {
	return r.root
}

// NewRouter creates light weight composible router that implements the http.Handler interface
// You can add handlers for specific paths, match on key words or regular expressions.
func NewRouter() *Router {
	r := &Router{}
	r.root = newRoot(r)
	return r
}

// Handler adds a handler function to the router for a path and method combination
func (r *Router) Handle(path []string, method string, handler http.HandlerFunc, mw ...Middleware) {

	r.root.insertMiddleware(path, mw)
	err := r.root.insertHandler(path, method, handler)
	if err != nil {
		panic(fmt.Sprintf("Failed to insert handler for path `/%s`: ", strings.Join(path, ",")) + err.Error())
	}
}

// Get add a GET method handler for the speficied path
func (r *Router) Get(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodGet, handler, mw...)
}

// Patch add a PATCH method handler for the speficied path
func (r *Router) Patch(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodPatch, handler, mw...)
}

// Put add a PUT method handler for the speficied path
func (r *Router) Put(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodPut, handler, mw...)
}

// Post add a POST method handler for the speficied path
func (r *Router) Post(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodPost, handler, mw...)
}

// Delete add a DELETE method handler for the speficied path
func (r *Router) Delete(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodDelete, handler, mw...)
}

// Head adds a HEAD method handler for the speficied path
func (r *Router) Head(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodHead, handler, mw...)
}

// // Options adds an OPTIONS method handler for the speficied path
// func (r *Router) Options(path string, handler http.HandlerFunc, mw ...Middleware) {
// 	r.Handle(r.splitPaths(path), http.MethodOptions, handler, mw...)
// }

func (r *Router) Static(path string, handler http.Handler, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodGet, handler.ServeHTTP, mw...)
}

// Connect adds a CONNECT method handler for the speficied path
func (r *Router) Connect(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodConnect, handler, mw...)
}

// Trace adds a TRACE method handler for the speficied path
func (r *Router) Trace(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.Handle(r.splitPaths(path), http.MethodTrace, handler, mw...)
}

// Use adds 1 or more middleware to the router.
// The middleware function will be applied to all handlers associated with this router
func (r *Router) Use(middleware ...Middleware) {
	r.root.insertMiddleware([]string{}, middleware)
}

func (r *Router) splitPaths(path string) []string {
	path = strings.TrimPrefix(r.root.Name+path, "/")
	paths := strings.Split(path, "/")
	for len(paths) > 0 && paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}
	return paths
}

// Router adds a set of routes to the router at a specified path.
// This is convinient for paths that share a common prefix reduing the need to specify the entire path for every handler.
func (router *Router) WithRoutes(path string, setRoutes func(r *Router)) {
	path = strings.TrimPrefix(path, "/")
	newRouter := NewRouter()
	setRoutes(newRouter)
	router.Mount(path, newRouter)
}

// NotFound replaces the routers default NotFound handler function.
func (r *Router) NotFound(handler http.HandlerFunc) {
	r.root.NotFound = handler
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

// Mount converts a top-level router into a top-level node in the parent router with a specified prefix
func (parent *Router) Mount(prefix string, subrouter *Router) {
	parent.root.insertNodeAt(prefix, subrouter.root)
}

// ServeHTTP is the Router's implementation of the http.Handler interface
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	router.root.ServeHTTP(w, r)
}
