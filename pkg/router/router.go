package router

import (
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Chi *chi.Mux
}

type HandleCtx func(ctx Ctx) error

func New() *Router {
	return &Router{
		Chi: chi.NewRouter(),
	}
}

func (rt *Router) handlerCtx(handler HandleCtx, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ResponseWriter(ctx) == nil {
		ctx = WithResponseWriter(ctx, w)
	}

	if Request(ctx) == nil {
		ctx = WithRequest(ctx, r)
	}

	r = r.WithContext(ctx)

	if err := handler(CtxRouter(r.Context())); err != nil {
		logging.Logging(ctx).Error(err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Router) Get(path string, handler HandleCtx) {
	rt.Chi.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) Post(path string, handler HandleCtx) {
	rt.Chi.Post(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) NotFound(handler HandleCtx) {
	rt.Chi.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) MethodNotAllowed(handler HandleCtx) {
	rt.Chi.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.Chi.Use(middlewares...)
}

func (rt *Router) Group(fn func(r *Router)) Router {
	im := rt.With()

	if fn != nil {
		fn(&im)
	}

	return im
}

func (rt *Router) With() Router {
	return Router{
		Chi: rt.Chi.With().(*chi.Mux),
	}
}

func (rt *Router) Route(pattern string, fn func(r *Router)) *chi.Mux {
	subRouter := New()

	fn(subRouter)
	rt.Chi.Mount(pattern, subRouter.Chi)

	return subRouter.Chi
}

func (rt *Router) ServeHTTP() http.HandlerFunc {
	return rt.Chi.ServeHTTP
}
