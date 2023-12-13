package gallery

import (
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func New(ctx router.Ctx) error {
	return tpl.Render(ctx, "gallery/new", tpl.Data{
		Data: ctx.Request.FormValue("title"),
	})
}
