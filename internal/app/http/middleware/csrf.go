package middleware

import (
	"github.com/gorilla/csrf"
	"net/http"
	"os"
	"strconv"
)

func Csrf(next http.Handler) http.Handler {
	csrfSecure, err := strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	if err != nil {
		csrfSecure = true
	}

	middleware := csrf.Protect(
		[]byte(os.Getenv("CSRF_KEY")),
		csrf.Secure(csrfSecure),
	)

	return middleware(next)
}
