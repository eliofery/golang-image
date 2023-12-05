package router

import (
	"context"
	"net/http"
)

const responseKey key = "response"

func WithResponseWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return context.WithValue(ctx, responseKey, w)
}

func ResponseWriter(ctx context.Context) http.ResponseWriter {
	val := ctx.Value(responseKey)

	resp, ok := val.(http.ResponseWriter)
	if !ok {
		return nil
	}

	return resp
}
