package cookie

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrCookieNotFound = errors.New("cookie не найдено")
)

const (
	Session = "session"
)

func New(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
}

func Set(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, New(name, value))
}

func Get(r *http.Request, name string) (string, error) {
	op := "cookie.Get"

	cookie, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrCookieNotFound)
	}

	return cookie.Value, nil
}

func Delete(w http.ResponseWriter, name string) {
	cookie := New(name, "")
	cookie.MaxAge = -1

	http.SetCookie(w, cookie)
}
