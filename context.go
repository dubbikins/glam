package glam

import (
	"fmt"
	"net/http"

	"github.com/dubbikins/glam/context"
)

const GlamContextKey context.ContextKey = "_!_glam_router_!_"

type glamContext struct {
	URLParams map[string]string
}

func newGlamContext() *glamContext {
	return &glamContext{
		URLParams: make(map[string]string),
	}
}

func GetURLParam(r *http.Request, key string) (string, bool) {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if ok {
		value, in := gctx.URLParams[fmt.Sprintf("{%s}", key)]
		return value, in
	}
	return "", false
}
func GetRegexURLParam(r *http.Request, key string) (string, bool) {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if ok {
		value, in := gctx.URLParams[fmt.Sprintf("`%s`", key)]
		return value, in
	}
	return "", false
}
func GetGlamContext(r *http.Request) *glamContext {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if ok {
		return gctx
	}
	return nil
}

func requestWithURLParam(r *http.Request, key, value string) *http.Request {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if !ok {
		gctx = newGlamContext()
	}

	gctx.URLParams[key] = value
	return ctx.Update(r, gctx)
}
