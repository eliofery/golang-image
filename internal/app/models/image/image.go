package image

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/static"
	"path/filepath"
)

var (
	extensions = []string{".jpg", ".jpeg", ".png"}
)

type Image struct {
	Path string
}

type Service struct {
	ctx router.Ctx
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Images(galleryID uint) ([]Image, error) {
	op := "model.image.Images"

	globPattern := filepath.Join(static.GalleriesDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var images []Image
	for _, file := range allFiles {
		if static.HasExtension(file, extensions) {
			images = append(images, Image{
				Path: file,
			})
		}
	}

	return images, nil
}
