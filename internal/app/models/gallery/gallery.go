package gallery

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrTitleAlreadyExists = errors.New("заголовок уже существует")
	ErrGalleryNotFound    = errors.New("галерея не существует")
)

type Gallery struct {
	ID     uint   `validate:"omitempty"`
	UserID uint   `validate:"required,min=1"`
	Title  string `validate:"required,max=255"`
}

type Service struct {
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Create(gallery *Gallery) error {
	op := "model.gallery.Create"

	db, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	err := v.Var(gallery.Title, "required,max=255")
	if err != nil {
		return err
	}

	row := db.QueryRow(`
        INSERT INTO galleries (user_id, title)
        VALUES ($1, $2) RETURNING id;`,
		gallery.UserID, gallery.Title,
	)

	err = row.Scan(&gallery.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return errors.Public(err, ErrTitleAlreadyExists.Error())
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) ByID(gallery *Gallery) error {
	op := "model.gallery.ById"

	db, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	err := v.Var(gallery.ID, "required,min=1")
	if err != nil {
		return err
	}

	row := db.QueryRow(`
        SELECT title, user_id
        FROM galleries WHERE id = $1;`,
		gallery.ID,
	)
	err = row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, ErrGalleryNotFound.Error())
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
