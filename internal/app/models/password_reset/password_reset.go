package pwreset

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/email"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Minute
)

var (
	ErrNotFount     = errors.New("заданный токен не существует")
	ErrTokenExpired = errors.New("токен просрочен")
)

type PasswordReset struct {
	ID        uint      `validate:"omitempty"`
	UserId    uint      `validate:"required,min=1"`
	TokenHash string    `validate:"required"`
	ExpiresAt time.Time `validate:"required,datetime"`
}

type Service struct {
	ctx router.Ctx

	email *email.Service
	user  *user.Service
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx: ctx,

		email: email.NewService(),
		user:  user.NewService(ctx),
	}
}

func (s *Service) Create(mail string) (*user.User, error) {
	op := "model.pwreset.Create"

	us := &user.User{
		Email: strings.ToLower(mail),
	}

	err := s.ctx.Validate.Var(us.Email, "required,email")
	if err != nil {
		return nil, err
	}

	row := s.ctx.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, us.Email)
	err = row.Scan(&us.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Public(err, user.ErrLoginOrPassword.Error())
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := rand.SessionToken()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pwReset := &PasswordReset{
		UserId:    us.ID,
		TokenHash: rand.HashToken(token),
		ExpiresAt: time.Now().Add(DefaultResetDuration),
	}

	row = s.ctx.DB.QueryRow(`
        INSERT INTO password_reset (user_id, token_hash, expires_at) VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO
        UPDATE SET token_hash = $2, expires_at = $3
        RETURNING id;`,
		us.ID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	vals := url.Values{"token": {token}}
	resetUrl := "http://localhost:8080/reset-pw?" + vals.Encode()

	err = s.email.Send(email.Email{
		From:    os.Getenv("EMAIL_SUPPORT"),
		To:      us.Email,
		Subject: "Восстановление пароля",
		Plaintext: `
           Вы запросили восстановление пароля.

           Если это были не вы проигнорируйте данное письмо.
           В противном случае перейдите по ссылке: ` + resetUrl,
		HTML: `
	      <h1>Вы запросили восстановление пароля.</h1>

	      <p>Если это были не вы проигнорируйте данное письмо.</p>
	      <p>В противном случае перейдите по ссылке: <a href="` + resetUrl + `">` + resetUrl + `</a></p>
	  `,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return us, nil
}

func (s *Service) Consume(password, token string) (string, error) {
	op := "model.pwreset.Consume"

	err := s.ctx.Validate.Var(password, "required,gte=10,lte=32")
	if err != nil {
		return token, err
	}

	err = s.ctx.Validate.Var(token, "required")
	if err != nil {
		return token, err
	}

	us := &user.User{
		Password: password,
	}

	pwReset := &PasswordReset{
		TokenHash: rand.HashToken(token),
		ExpiresAt: time.Now().Add(DefaultResetDuration),
	}

	row := s.ctx.DB.QueryRow(`
        SELECT password_reset.id, password_reset.expires_at, users.id, users.email
        FROM password_reset
        INNER JOIN users ON users.id = password_reset.user_id
        WHERE password_reset.token_hash = $1;`, pwReset.TokenHash)
	err = row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &us.ID, &us.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return token, errors.Public(err, ErrNotFount.Error())
		}
		return token, fmt.Errorf("%s: %w", op, err)
	}

	err = s.Delete(pwReset)
	if err != nil {
		return token, fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return token, errors.Public(nil, fmt.Sprintf("%s: %s", ErrTokenExpired, token))
	}

	err = s.user.UpdatePassword(us)
	if err != nil {
		return token, fmt.Errorf("%s: %w", op, err)
	}

	err = session.NewService(s.ctx).Create(&session.Session{UserID: us.ID})
	if err != nil {
		return token, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Service) Delete(pwReset *PasswordReset) error {
	op := "model.pwreset.Delete"

	_, err := s.ctx.DB.Exec(`
        DELETE FROM password_reset
        WHERE id = $1;`, pwReset.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
