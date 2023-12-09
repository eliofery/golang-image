package notfound

import (
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Page404(ctx router.Ctx) error {
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

	return tpl.Render(ctx, "error/404", tpl.Data{})
}

func Page405(ctx router.Ctx) error {
	ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

	return tpl.Render(ctx, "error/404", tpl.Data{})
}
