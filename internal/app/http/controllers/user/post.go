package user

import (
	pwreset "github.com/eliofery/golang-image/internal/app/models/password_reset"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Create(ctx router.Ctx) error {
	userData := user.User{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	service := user.NewService(ctx)
	err := service.SignUp(&userData)
	if err != nil {
		ctx.Info(err.Error())
		ctx.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signup", tpl.Data{
			Data:   userData,
			Errors: []error{err},
		})
	}

	http.Redirect(ctx.ResponseWriter, ctx.Request, "/user", http.StatusFound)

	return nil
}

func Auth(ctx router.Ctx) error {
	userData := user.User{
		Email:    ctx.FormValue("email"),
		Password: ctx.FormValue("password"),
	}

	service := user.NewService(ctx)
	err := service.SignIn(&userData)
	if err != nil {
		ctx.Info(err.Error())
		ctx.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signin", tpl.Data{
			Data:   userData,
			Errors: []error{err},
		})
	}

	http.Redirect(ctx.ResponseWriter, ctx.Request, "/user", http.StatusFound)

	return nil
}

func Logout(ctx router.Ctx) error {
	cookie.Delete(ctx.ResponseWriter, cookie.Session)

	http.Redirect(ctx.ResponseWriter, router.Request(ctx), "/signin", http.StatusFound)

	return nil
}

func ProcessForgotPassword(ctx router.Ctx) error {
	userData := user.User{
		Email: ctx.FormValue("email"),
	}

	service := pwreset.NewService(ctx)
	err := service.Create(&userData)
	if err != nil {
		ctx.Info(err.Error())
		ctx.WriteHeader(http.StatusInternalServerError)

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
	data := &struct {
		Password string
		Token    string
	}{
		Password: ctx.FormValue("password"),
		Token:    ctx.FormValue("token"),
	}

	service := pwreset.NewService(ctx)
	err := service.Consume(data)
	if err != nil {
		var pubErr errors.PublicError
		if errors.As(err, &pubErr) {
			ctx.Info(pubErr.Public())
		} else {
			ctx.Info(err.Error())
		}

		ctx.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/reset-pw", tpl.Data{
			Data:   data,
			Errors: []error{err},
		})
	}

	http.Redirect(ctx.ResponseWriter, ctx.Request, "/user", http.StatusFound)

	return nil
}
