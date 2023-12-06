package gallery

import (
	"context"
	"fmt"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrTitleAlreadyExists = errors.New("заголовок уже существует")
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
