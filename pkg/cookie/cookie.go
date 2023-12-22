package cookie

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/eliofery/golang-image/pkg/router"
	"net/http"
)

var (
	ErrCookieNotFound = errors.New("cookie не найдено")
)

const (
	Session = "session"
	Message = "message"
)

func New(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    string([]rune(value)),
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

func SetMessage(w http.ResponseWriter, message string) {
	Set(w, Message, base64.StdEncoding.EncodeToString([]byte(message)))
}

func GetMessage(ctx router.Ctx) string {
	op := "cookie.GetMessage"

	messageEncode, err := Get(ctx.Request, Message)
	if err != nil {
		ctx.Logger.Info(fmt.Sprintf("%s: %s", op, err))
	}

	message, err := base64.StdEncoding.DecodeString(messageEncode)
	if err != nil {
		ctx.Logger.Info(fmt.Sprintf("%s: %s", op, err))
		return ""
	}

	return string(message)
}
