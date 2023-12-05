package router

import (
	"context"
	"net/http"
)

type key string

const requestKey key = "request"

func WithRequest(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, r)
}

func Request(ctx context.Context) *http.Request {
	val := ctx.Value(requestKey)

	request, ok := val.(*http.Request)
	if !ok {
		return nil
	}

	return request
}
