package gallery

import (
	"github.com/eliofery/golang-image/internal/app/models/gallery"
	"github.com/eliofery/golang-image/internal/app/models/image"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/go-chi/chi/v5"
	"net/http"
	"path/filepath"
	"strconv"
)

var (
	errNotAllowed = errors.New("нет доступа к галереи")
)

func Index(ctx router.Ctx) error {
	sGallery := gallery.NewService(ctx)

	userData := user.CtxUser(ctx)
	galleriesData, err := sGallery.ByUserID(userData.ID)
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	message, err := cookie.GetMessage(ctx.Request)
	if err != nil {
		ctx.Logger.Info(err.Error())
	}
	cookie.Delete(ctx.ResponseWriter, cookie.Message)

	return tpl.Render(ctx, "gallery/index", tpl.Data{
		Data:     galleriesData,
		Messages: []any{message},
	})
}

func Show(ctx router.Ctx) error {
	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	sGallery := gallery.NewService(ctx)
	sImage := image.NewService(ctx)

	galleryData, err := sGallery.ByID(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	userData := user.CtxUser(ctx)
	if galleryData.UserID != userData.ID {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return tpl.Render(ctx, "error/405", tpl.Data{})
	}

	images, err := sImage.Images(galleryData.ID)
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/500", tpl.Data{})
	}

	data := struct {
		ID     uint
		Title  string
		Images []image.Image
	}{
		ID:     galleryData.ID,
		Title:  galleryData.Title,
		Images: images,
	}

	return tpl.Render(ctx, "gallery/show", tpl.Data{
		Data: data,
	})
}

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

	sGallery := gallery.NewService(ctx)
	sImage := image.NewService(ctx)

	galleryData, err := sGallery.ByID(uint(id))
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

	images, err := sImage.Images(galleryData.ID)
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/500", tpl.Data{})
	}

	data := struct {
		ID     uint
		Title  string
		Images []image.Image
	}{
		ID:     galleryData.ID,
		Title:  galleryData.Title,
		Images: images,
	}

	return tpl.Render(ctx, "gallery/edit", tpl.Data{
		Data:     data,
		Messages: []any{message},
	})
}

func Image(ctx router.Ctx) error {
	fileName := filepath.Base(chi.URLParam(ctx.Request, "filename"))

	galleryID, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())

		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	sImage := image.NewService(ctx)
	imageData, err := sImage.Image(uint(galleryID), fileName)
	if err != nil {
		ctx.Logger.Info(err.Error())

		if errors.Is(err, image.ErrNotFound) {
			ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
			return nil
		}

		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	http.ServeFile(ctx.ResponseWriter, ctx.Request, imageData.FilePath)

	return nil
}
