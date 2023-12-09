package session

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Session struct {
	ID        uint   `validate:"omitempty"`
	UserID    uint   `validate:"required,min=1"`
	TokenHash string `validate:"required"`
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

func (s *Service) Create(session *Session) error {
	op := "model.session.SignUp"

	token, err := rand.SessionToken()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	session.TokenHash = rand.HashToken(token)

	err = s.validate.Struct(session)
	if err != nil {
		return err
	}

	row := s.db.QueryRow(`
        INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)
        ON CONFLICT (user_id) DO
        UPDATE SET token_hash = $2
        RETURNING id;`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cookie.Set(s.writer, cookie.Session, token)

	return nil
}
