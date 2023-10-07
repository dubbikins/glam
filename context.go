package glam

import (
	"net/http"

	"context"
)

type ContextKey string

const routerContextKey ContextKey = "_!_router_!_"

func Query(r *http.Request, key string) (string, bool) {
	return r.URL.Query().Get(key), r.URL.Query().Has(key)
}

func withParam(r *http.Request, key, value string) *http.Request {
	ctx, ok := r.Context().Value(routerContextKey).(map[string]string)
	if !ok {
		ctx = map[string]string{
			key[1 : len(key)-1]: value,
		}
		return r.WithContext(context.WithValue(r.Context(), routerContextKey, ctx))
	}
	if _, in := ctx[key[1:len(key)-1]]; in {
		logger.Warn("duplicate param key identified")
	}
	ctx[key[1:len(key)-1]] = value
	return r
}

func withRegex(r *http.Request, key, value string) *http.Request {
	ctx, ok := r.Context().Value(routerContextKey).(map[string]string)
	start, sep := getRegexKeyIndices(key)
	if !ok {
		ctx = map[string]string{
			key[start:sep]: value,
		}
		return r.WithContext(context.WithValue(r.Context(), routerContextKey, ctx))
	}
	if _, in := ctx[key[start:sep]]; in {
		logger.Warn("duplicate regex key identified")
	}
	ctx[key[start:sep]] = value
	return r
}

func GetParam(r *http.Request, key string) (string, bool) {
	ctx, ok := r.Context().Value(routerContextKey).(map[string]string)
	if ok {
		value, in := ctx[key]
		return value, in
	}
	return "", false
}
