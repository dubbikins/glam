package context

import (
	"context"
	"net/http"
)

type ContextKey string

type Context[T any] struct {
	parent context.Context
	key    ContextKey
}

func NewContext[T any](r *http.Request, key ContextKey) *Context[T] {
	ctx := &Context[T]{
		parent: r.Context(),
	}
	return ctx
}
func (ctx *Context[T]) Update(r *http.Request, value T) *http.Request {
	ctx.parent = context.WithValue(ctx.parent, ctx.key, value)
	return r.WithContext(ctx.parent)
}

func (ctx *Context[T]) Get() (T, bool) {
	value, ok := ctx.parent.Value(ctx.key).(T)
	return value, ok
}
