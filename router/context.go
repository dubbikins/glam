package router

import (
	"context"
	"net/http"
)

const RouterContextKey = "_!_glam_router_!_"

type Context struct {
	URLParams map[string] string `json:"urlParams"`
}

func NewRouterContext(urlParams map[string]string) *Context {
	return &Context{
		URLParams: urlParams,
	}
}
func RouterContext(ctx context.Context) *Context {
	return ctx.Value(RouterContextKey).(*Context)
}
func (rctx *Context) InsertURLParams(ctx context.Context) context.Context {
	return context.WithValue(ctx, RouterContextKey, rctx)
}
func GetURLParam(r *http.Request, key string) string {
	if ctx := RouterContext(r.Context()); ctx != nil {
		return ctx.URLParam(key)
	}
	return ""
}


func (r *Context) URLParam(key string) string {
	if value, in := r.URLParams[key]; in {
		return value
	}
	return ""
}