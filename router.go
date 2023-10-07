package glam

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/dubbikins/glam/logging"
)

var logger = logging.Logger

type Router struct {
	Name           string                  //16
	ParamChild     *Router                 //8
	Children       children                //8
	RegexpChildren children                //8
	StaticChildren children                //8
	Handlers       map[string]http.Handler //8
	Middleware     []Middleware            //24
	notFound       http.Handler
	regexMatcher   *regexp.Regexp
	parent         *Router
}

// NewRouter creates light weight composible router that implements the http.Handler interface
// You can add handlers for specific paths, match on key words or regular expressions.
func NewRouter() *Router {
	r := &Router{
		Name:           "",
		Children:       newChildren(),
		StaticChildren: newChildren(),
		RegexpChildren: newChildren(),
	}
	return r
}

// Handler adds a handler function to the router for a path and method combination
func (r *Router) handle(path []string, method string, handler http.HandlerFunc, mw ...Middleware) {

	r.insertMiddleware(path, mw)
	err := r.insertHandler(path, method, handler)
	if err != nil {
		panic(fmt.Sprintf("Failed to insert handler for path `/%s`: ", strings.Join(path, ",")) + err.Error())
	}
}

// Get add a GET method handler for the speficied path
func (r *Router) Get(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodGet, handler, mw...)
}

// Patch add a PATCH method handler for the speficied path
func (r *Router) Patch(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodPatch, handler, mw...)
}

// Put add a PUT method handler for the speficied path
func (r *Router) Put(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodPut, handler, mw...)
}

// Post add a POST method handler for the speficied path
func (r *Router) Post(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodPost, handler, mw...)
}

// Delete add a DELETE method handler for the speficied path
func (r *Router) Delete(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodDelete, handler, mw...)
}

// Head adds a HEAD method handler for the speficied path
func (r *Router) Head(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodHead, handler, mw...)
}

// // Options adds an OPTIONS method handler for the speficied path
// func (r *Router) Options(path string, handler http.HandlerFunc, mw ...Middleware) {
// 	r.handle(r.splitPaths(path), http.MethodOptions, handler, mw...)
// }

// Connect adds a CONNECT method handler for the speficied path
func (r *Router) Connect(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodConnect, handler, mw...)
}

// Trace adds a TRACE method handler for the speficied path
func (r *Router) Trace(path string, handler http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodTrace, handler, mw...)
}

// Use adds 1 or more middleware to the router.
// The middleware function will be applied to all handlers associated with this router
func (r *Router) Use(middleware ...Middleware) {
	r.insertMiddleware([]string{}, middleware)
}

func (r *Router) splitPaths(path string) []string {
	path = strings.TrimPrefix(r.Name+path, "/")
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
	r.notFound = handler
}

var defaultNotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

// Mount converts a top-level router into a top-level node in the parent router with a specified prefix
func (parent *Router) Mount(prefix string, subrouter *Router) {
	subrouter.Name = strings.TrimPrefix(prefix, "/")
	subrouter.parent = parent
	parent.insertNodeAt(prefix, subrouter)
}

func (r *Router) Static(path string, handler http.Handler, mw ...Middleware) {
	r.handle(r.splitPaths(path+StaticIdentifier), http.MethodGet, handler.ServeHTTP, mw...)
}

func (r *Router) Websocket(path string, ws http.HandlerFunc, mw ...Middleware) {
	r.handle(r.splitPaths(path), http.MethodGet, ws, mw...)
}
