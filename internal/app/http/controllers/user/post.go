package user

import (
	pwreset "github.com/eliofery/golang-image/internal/app/models/password_reset"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Create(ctx router.Ctx) error {
	w, r, l := router.ResponseWriter(ctx), router.Request(ctx), logging.Logging(ctx)

	userData := user.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	service := user.NewService(ctx)
	err := service.SignUp(&userData)
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signup", tpl.Data{
			Data:   userData,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}

func Auth(ctx router.Ctx) error {
	w, r, l := router.ResponseWriter(ctx), router.Request(ctx), logging.Logging(ctx)

	userData := user.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	service := user.NewService(ctx)
	err := service.SignIn(&userData)
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signin", tpl.Data{
			Data:   userData,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}

func Logout(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)

	cookie.Delete(w, cookie.Session)

	http.Redirect(w, router.Request(ctx), "/signin", http.StatusFound)

	return nil
}

func ProcessForgotPassword(ctx router.Ctx) error {
	r, w, l := router.Request(ctx), router.ResponseWriter(ctx), logging.Logging(ctx)

	userData := user.User{
		Email: r.FormValue("email"),
	}

	service := pwreset.NewService(ctx)
	err := service.Create(&userData)
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/forgot-pw", tpl.Data{
			Data:   userData,
			Errors: []error{err},
		})
	}

	return tpl.Render(ctx, "user/check-email", tpl.Data{
		Data: userData,
	})
}

func ProcessResetPassword(ctx router.Ctx) error {
	r, w, l := router.Request(ctx), router.ResponseWriter(ctx), logging.Logging(ctx)

	data := &struct {
		Password string
		Token    string
	}{
		Password: r.FormValue("password"),
		Token:    r.FormValue("token"),
	}

	service := pwreset.NewService(ctx)
	err := service.Consume(data)
	if err != nil {
		var pubErr errors.PublicError
		if errors.As(err, &pubErr) {
			l.Info(pubErr.Public())
		} else {
			l.Info(err.Error())
		}

		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/reset-pw", tpl.Data{
			Data:   data,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}
