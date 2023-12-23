package gallery

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/image"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
)

var (
	ErrGalleryAlreadyExists = errors.New("галерея уже существует")
	ErrGalleryNotFound      = errors.New("галерея не найдена")
)

type Gallery struct {
	ID     uint   `validate:"omitempty"`
	UserID uint   `validate:"required,min=1"`
	Title  string `validate:"required,max=255"`
}

type Service struct {
	ctx router.Ctx

	user *user.User
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx: ctx,

		user: user.CtxUser(ctx),
	}
}

func (s *Service) Create(title string) (*Gallery, error) {
	op := "model.gallery.Create"

	gallery := &Gallery{
		UserID: s.user.ID,
		Title:  title,
	}

	err := s.ctx.Validate.Struct(gallery)
	if err != nil {
		return nil, err
	}

	row := s.ctx.DB.QueryRow(`
        INSERT INTO galleries (user_id, title)
        VALUES ($1, $2) RETURNING id;`,
		gallery.UserID, gallery.Title,
	)
	err = row.Scan(&gallery.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return nil, errors.Public(err, ErrGalleryAlreadyExists.Error())
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return gallery, nil
}

func (s *Service) ByID(id uint) (*Gallery, error) {
	op := "model.gallery.ById"

	gallery := &Gallery{
		ID: id,
	}

	err := s.ctx.Validate.Var(gallery.ID, "required,min=1")
	if err != nil {
		return nil, err
	}

	row := s.ctx.DB.QueryRow(`
        SELECT title, user_id
        FROM galleries WHERE id = $1;`,
		gallery.ID,
	)
	err = row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Public(err, ErrGalleryNotFound.Error())
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return gallery, nil
}

func (s *Service) ByUserID(id uint) ([]Gallery, error) {
	op := "model.gallery.ById"

	err := s.ctx.Validate.Var(id, "required,min=1")
	if err != nil {
		return nil, err
	}

	rows, err := s.ctx.DB.Query(`
        SELECT id, user_id, title
        FROM galleries
        WHERE user_id = $1;`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var galleries []Gallery
	for rows.Next() {
		var gallery Gallery

		err = rows.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return galleries, nil
}

func (s *Service) Update(gallery *Gallery) error {
	op := "model.gallery.Delete"

	_, err := s.ctx.DB.Exec(`UPDATE galleries SET title = $1 WHERE id = $2;`, gallery.Title, gallery.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return errors.Public(err, ErrGalleryAlreadyExists.Error())
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Delete(id uint) error {
	op := "model.gallery.Delete"

	_, err := s.ctx.DB.Exec(`DELETE FROM galleries WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = os.RemoveAll(image.GalleriesDir(id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
