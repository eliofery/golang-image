package user

import (
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/index", tpl.Data{})
}

func SignUp(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signup", tpl.Data{})
}

func SignIn(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signin", tpl.Data{})
}

func ForgotPassword(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/forgot-pw", tpl.Data{})
}

func ResetPassword(ctx router.Ctx) error {
	r := router.Request(ctx)

	data := &struct {
		Token string
	}{
		Token: r.FormValue("token"),
	}

	return tpl.Render(ctx, "user/reset-pw", tpl.Data{
		Data: data,
	})
}
