package middleware

import (
	"github.com/eliofery/golang-image/pkg/cookie"
	"net/http"
)

func RemoveMessage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie.Delete(w, cookie.Message)

		next.ServeHTTP(w, r)
	})
}
