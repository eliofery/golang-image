package gallery

import (
	"github.com/eliofery/golang-image/internal/app/models/gallery"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var (
	errNotAllowed = errors.New("нет доступа к галереи")
)

func New(ctx router.Ctx) error {
	message, err := cookie.GetMessage(ctx.Request)
	if err != nil {
		ctx.Logger.Info(err.Error())
	}
	cookie.Delete(ctx.ResponseWriter, cookie.Message)

	return tpl.Render(ctx, "gallery/new", tpl.Data{
		Data:     ctx.Request.FormValue("title"),
		Messages: []any{message},
	})
}

func Edit(ctx router.Ctx) error {
	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{
			//Errors: []error{err},
		})
	}

	service := gallery.NewService(ctx)

	galleryData, err := service.ByID(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{
			//Errors: []error{err},
		})
	}

	userData := user.CtxUser(ctx)
	if galleryData.UserID != userData.ID {
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return tpl.Render(ctx, "error/405", tpl.Data{
			Errors: []error{errors.Public(err, errNotAllowed.Error())},
		})
	}

	message, err := cookie.GetMessage(ctx.Request)
	if err != nil {
		ctx.Logger.Info(err.Error())
	}
	cookie.Delete(ctx.ResponseWriter, cookie.Message)

	return tpl.Render(ctx, "gallery/edit", tpl.Data{
		Data:     galleryData,
		Messages: []any{message},
	})
}
