package image

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/static"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

var (
	ErrNotFound = errors.New("изображение не найдено")

	extensions   = []string{".jpg", ".jpeg", ".png"}
	contentTypes = []string{"image/jpeg", "image/jpg", "image/png"}
)

type Image struct {
	GalleryID uint
	FilePath  string
	FileName  string
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

	globPattern := filepath.Join(GalleriesDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var images []Image
	for _, file := range allFiles {
		if static.HasExtension(file, extensions) {
			images = append(images, Image{
				GalleryID: galleryID,
				FilePath:  file,
				FileName:  url.PathEscape(filepath.Base(file)),
			})
		}
	}

	return images, nil
}

func (s *Service) Image(galleryID uint, filename string) (Image, error) {
	op := "model.image.Image"

	imagePath := filepath.Join(GalleriesDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Image{}, ErrNotFound
		}

		return Image{}, fmt.Errorf("%s: %w", op, err)
	}

	return Image{
		GalleryID: galleryID,
		FilePath:  imagePath,
		FileName:  url.PathEscape(filepath.Base(filename)),
	}, nil
}

func (s *Service) Delete(galleryID uint, filename string) error {
	op := "model.image.Delete"

	image, err := s.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = os.Remove(image.FilePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) CreateImage(galleryID uint, filename string, contents io.ReadSeeker) error {
	op := "model.gallery.CreateImage"

	err := errors.CheckContentType(contents, contentTypes)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = errors.CheckExtension(filename, extensions)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	galleryDir := GalleriesDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	imagePath := filepath.Join(galleryDir, filename)
	file, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	_, err = io.Copy(file, contents)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func GalleriesDir(id uint) string {
	dir := os.Getenv("STATIC_DIR")
	if dir == "" {
		dir = "internal/static"
	}

	return filepath.Join(dir, fmt.Sprintf("galleries/%d", id))
}
