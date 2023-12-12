package gallery

import (
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	return tpl.Render(ctx, "gallery/index", tpl.Data{
		Data: ctx.Request.FormValue("title"),
	})
}
