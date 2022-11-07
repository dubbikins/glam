package glam

import (
	"fmt"
	"net/http"

	"github.com/dubbikins/glam/context"
)

const GlamContextKey context.ContextKey = "_!_glam_router_!_"

type glamContext struct {
	Params map[string]string
}

func newGlamContext() *glamContext {
	return &glamContext{
		Params: make(map[string]string),
	}
}

func GetParam(r *http.Request, key string) (string, bool) {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if ok {
		value, in := gctx.Params[fmt.Sprintf("{%s}", key)]
		return value, in
	}
	return "", false
}
func GetRegexParam(r *http.Request, key string) (string, bool) {
	ctx := context.NewContext[*glamContext](r, GlamContextKey)
	gctx, ok := ctx.Get()
	if ok {
		value, in := gctx.Params[fmt.Sprintf("(%s)", key)]
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

	gctx.Params[key] = value
	return ctx.Update(r, gctx)
}
