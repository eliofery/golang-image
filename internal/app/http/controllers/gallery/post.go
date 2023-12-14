package gallery

import (
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/gallery"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Create(ctx router.Ctx) error {
	service := gallery.NewService(ctx)

	galleryData, err := service.Create(ctx.Request.FormValue("title"))
	if err != nil {
		ctx.Info(err.Error())
		ctx.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "gallery/new", tpl.Data{
			Errors: []error{err},
		})
	}

	editPath := fmt.Sprintf("/gallery/%d/edit", galleryData.ID)

	http.Redirect(ctx.ResponseWriter, ctx.Request, editPath, http.StatusFound)

	return nil
}
