package gallery

import (
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/gallery"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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

func Update(ctx router.Ctx) error {
	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	service := gallery.NewService(ctx)

	galleryData, err := service.ByID(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	userData := user.CtxUser(ctx)
	if galleryData.UserID != userData.ID {
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return tpl.Render(ctx, "error/405", tpl.Data{
			Errors: []error{errors.Public(err, errNotAllowed.Error())},
		})
	}

	galleryData.Title = ctx.Request.FormValue("title")

	err = service.Update(galleryData)
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/500", tpl.Data{
			Errors: []error{err},
		})
	}

	return tpl.Render(ctx, "gallery/edit", tpl.Data{
		Data:     galleryData,
		Messages: []any{"Галерея успешно обновлена"},
	})
}
