package session

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
)

type Session struct {
	ID        uint   `validate:"omitempty"`
	UserID    uint   `validate:"required,min=1"`
	TokenHash string `validate:"required"`
}

type Service struct {
	ctx router.Ctx
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Create(session *Session) error {
	op := "model.session.SignUp"

	token, err := rand.SessionToken()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	session.TokenHash = rand.HashToken(token)

	err = s.ctx.Validate.Struct(session)
	if err != nil {
		return err
	}

	row := s.ctx.DB.QueryRow(`
        INSERT INTO sessions (user_id, token_hash) VALUES ($1, $2)
        ON CONFLICT (user_id) DO
        UPDATE SET token_hash = $2
        RETURNING id;`,
		session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cookie.Set(s.ctx.ResponseWriter, cookie.Session, token)

	return nil
}
