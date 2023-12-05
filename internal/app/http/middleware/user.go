package middleware

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"net/http"
)

func SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := user.GetCurrentUser(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := user.WithUser(r.Context(), userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
