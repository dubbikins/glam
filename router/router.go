package router 

import (
	"fmt"
	"net/http"

)

type UrlParameter struct {

}

type Router struct {
	PathParameters map[string]string
	root *Node
	Parent *Router
	Result string

	middleware []Middleware
	handler http.Handler
	NotFoundHandler http.Handler
	MiddlewareHandler http.Handler
}

func NewRouter( ) *Router {
	return &Router{
		root: NewRoot(),
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			w.WriteHeader(http.StatusNotFound)
		}),
	}
}

func (c *Router ) PrintRouterTree(){
	//
}
func (c *Router) Handle( path, method string, handler http.HandlerFunc){
	c.root.InsertRouteHandler(c.root.Name+path, http.MethodGet, handler)
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
	for _, mw := range middleware {
		r.useAt( mw,r.root.Name)
	}
}
func (r *Router) useAt(middleware Middleware, path string ) {
	r.root.InsertRouteMiddleware(path, middleware)
}
func (r *Router) UseAt(middleware Middleware, paths ...string) {
	if len(paths) == 0 {
		r.useAt( middleware, r.root.Name,)
	} else {
		for _, path := range paths {
			r.useAt(middleware,r.root.Name+path)
		}
	}
}

func (router *Router) Router (path string, setRoutes func (r *Router)) {
	newRouter := NewRouter()
	setRoutes(newRouter)
	router.Mount(path,newRouter)
}

func (r *Router) NotFound(handler http.HandlerFunc) {
	r.NotFoundHandler = handler
}

func (r *Router) GetUrlParam (name string) string {
	for r.Parent != nil {
		r = r.Parent
	}
	return r.PathParameters[fmt.Sprintf("{%s}", name)]
}

func (parent *Router) Mount(path string, r *Router) {
	r.Parent = parent
	parent.root.InsertNodeAt(path, r.root)
}

func (r *Router) Handler(req *http.Request) (http.Handler, *http.Request) {
	middleware, handler, urlParams, err := r.root.MatchRequest(req)
	if len(urlParams ) > 0 {
		rctx := NewRouterContext(urlParams)
		ctx := rctx.InsertURLParams(req.Context())
		req = req.WithContext(ctx)
	}
	if handler == nil  || err != nil {
		handler = r.NotFoundHandler
	}
	if len(middleware) == 0 {
		return handler, req
	}
	h := middleware[len(middleware)-1](handler)
	for i := len(middleware) - 2; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h, req
}

func (router *Router) PrintTree() {
	router.root.PrintTree()
}
func (router *Router) String() {
	router.root.String()
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h, reqWithRouterContext := router.Handler(r)
	h.ServeHTTP(w, reqWithRouterContext)
}

