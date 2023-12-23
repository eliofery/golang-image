package gallery

import (
	"fmt"
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

	cookie.SetMessage(ctx.ResponseWriter, "Галерея создана")

	editPath := fmt.Sprintf("/gallery/%d/edit", galleryData.ID)
	http.Redirect(ctx.ResponseWriter, ctx.Request, editPath, http.StatusFound)

	return nil
}

func Update(ctx router.Ctx) error {
	service := gallery.NewService(ctx)

	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	galleryData, err := service.ByID(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	if galleryData.UserID != user.CtxUser(ctx).ID {
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

	cookie.SetMessage(ctx.ResponseWriter, "Галерея успешно обновлена")

	editPath := fmt.Sprintf("/gallery/%d/edit", galleryData.ID)
	http.Redirect(ctx.ResponseWriter, ctx.Request, editPath, http.StatusFound)

	return nil
}

func Delete(ctx router.Ctx) error {
	service := gallery.NewService(ctx)

	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	galleryData, err := service.ByID(uint(id))
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	if galleryData.UserID != user.CtxUser(ctx).ID {
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return tpl.Render(ctx, "error/405", tpl.Data{
			Errors: []error{errors.Public(err, errNotAllowed.Error())},
		})
	}

	err = service.Delete(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	cookie.SetMessage(ctx.ResponseWriter, "Галерея успешно удалена")

	http.Redirect(ctx.ResponseWriter, ctx.Request, "/gallery", http.StatusFound)

	return nil
}

func DeleteImage(ctx router.Ctx) error {
	sGallery := gallery.NewService(ctx)
	sImage := image.NewService(ctx)

	filename := filepath.Base(chi.URLParam(ctx.Request, "filename"))
	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	galleryData, err := sGallery.ByID(uint(id))
	if err != nil {
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

	err = sImage.Delete(galleryData.ID, filename)
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{})
	}

	cookie.SetMessage(ctx.ResponseWriter, "Изображение успешно удалено")

	editPath := fmt.Sprintf("/gallery/%d/edit", galleryData.ID)
	http.Redirect(ctx.ResponseWriter, ctx.Request, editPath, http.StatusFound)

	return nil
}

func UploadImage(ctx router.Ctx) error {
	sGallery := gallery.NewService(ctx)
	sImage := image.NewService(ctx)

	id, err := strconv.Atoi(chi.URLParam(ctx.Request, "id"))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{
			Errors: []error{err},
		})
	}

	galleryData, err := sGallery.ByID(uint(id))
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound)

		return tpl.Render(ctx, "error/404", tpl.Data{
			Errors: []error{err},
		})
	}

	if galleryData.UserID != user.CtxUser(ctx).ID {
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return tpl.Render(ctx, "error/405", tpl.Data{
			Errors: []error{errors.Public(err, errNotAllowed.Error())},
		})
	}

	err = ctx.Request.ParseMultipartForm(5 << 20) // 5 * 1024 * 1024 -> 5mb
	if err != nil {
		ctx.Logger.Info(err.Error())
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "error/500", tpl.Data{
			Errors: []error{err},
		})
	}

	fileHeaders := ctx.Request.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			ctx.Logger.Info(err.Error())
			ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

			return tpl.Render(ctx, "error/500", tpl.Data{
				Errors: []error{err},
			})
		}

		// TODO разобраться с корректным закрытием файла
		defer file.Close()

		err = sImage.CreateImage(galleryData.ID, fileHeader.Filename, file)
		if err != nil {
			var (
				fileErr    errors.FileError
				errMessage error
			)

			if errors.As(err, &fileErr) {
				errMessage = errors.Public(err, fmt.Sprintf("%v не верный тип файла. Допустимые типы jpg, jpeg, png", fileHeader.Filename))
			} else {
				errMessage = err
			}

			ctx.Logger.Info(errMessage.Error())

			ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)

			return tpl.Render(ctx, "error/405", tpl.Data{
				Errors: []error{errMessage},
			})
		}
	}

	cookie.SetMessage(ctx.ResponseWriter, "Изображения успешно загружены")

	editPath := fmt.Sprintf("/gallery/%d/edit", galleryData.ID)
	http.Redirect(ctx.ResponseWriter, ctx.Request, editPath, http.StatusFound)

	return nil
}
