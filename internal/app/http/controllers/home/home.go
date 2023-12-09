package home

import (
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	token, err := rand.SessionToken()
	if err != nil {
		ctx.Logger.Info("не удалось получить токен", err)
	}

	data := tpl.Data{
		Meta: tpl.Meta{
			Title: "Главная",
		},
		Data: struct {
			Token string
		}{
			Token: token,
		},
		Errors: tpl.PublicErrors(
			"ошибка 1",
			"ошибка 2",
			"ошибка 3",
		),
	}

	return tpl.Render(ctx, "home", data)
}
