package gallery

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
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
	ctx      router.Ctx
	writer   http.ResponseWriter
	db       *sql.DB
	validate *validator.Validate
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx:      ctx,
		writer:   ctx.ResponseWriter,
		db:       ctx.DB,
		validate: ctx.Validate,
	}
}

func (s *Service) Create(gallery *Gallery) error {
	op := "model.gallery.Create"

	err := s.validate.Var(gallery.Title, "required,max=255")
	if err != nil {
		return err
	}

	row := s.db.QueryRow(`
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

	err := s.validate.Var(gallery.ID, "required,min=1")
	if err != nil {
		return err
	}

	row := s.db.QueryRow(`
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

func (s *Service) ByUserID(us *user.User) ([]Gallery, error) {
	op := "model.gallery.ById"

	err := s.validate.Var(us.ID, "required,min=1")
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
        SELECT id, title
        FROM galleries
        WHERE user_id = $1;`,
		us.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var galleries []Gallery
	for rows.Next() {
		var gallery Gallery

		err = rows.Scan(&gallery.ID, &gallery.Title)
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
