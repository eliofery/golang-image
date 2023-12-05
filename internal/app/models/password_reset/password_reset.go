package pwreset

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/email"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/validate"
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
	ctx context.Context
}

func NewService(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
	}
}

func (s *Service) Create(us *user.User) error {
	op := "model.pwreset.Create"

	d, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	us.Email = strings.ToLower(us.Email)

	err := v.Var(us.Email, "required,email")
	if err != nil {
		return err
	}

	row := d.QueryRow(`SELECT id FROM users WHERE email = $1`, us.Email)
	err = row.Scan(&us.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, user.ErrLoginOrPassword.Error())
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	token, err := rand.SessionToken()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	pwReset := &PasswordReset{
		UserId:    us.ID,
		TokenHash: rand.HashToken(token),
		ExpiresAt: time.Now().Add(DefaultResetDuration),
	}

	row = d.QueryRow(`
        INSERT INTO password_reset (user_id, token_hash, expires_at) VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO
        UPDATE SET token_hash = $2, expires_at = $3
        RETURNING id;`,
		us.ID, pwReset.TokenHash, pwReset.ExpiresAt)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	vals := url.Values{"token": {token}}
	resetUrl := "http://localhost:8080/reset-pw?" + vals.Encode()

	emailService := email.NewService()
	err = emailService.Send(email.Email{
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
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Consume(data *struct{ Password, Token string }) error {
	op := "model.pwreset.Consume"

	d, v := database.CtxDatabase(s.ctx), validate.Validation(s.ctx)

	err := v.Var(data.Password, "required,gte=10,lte=32")
	if err != nil {
		return err
	}

	us := &user.User{}
	pwReset := &PasswordReset{
		ExpiresAt: time.Now().Add(DefaultResetDuration),
	}
	tokenHash := rand.HashToken(data.Token)

	row := d.QueryRow(`
        SELECT password_reset.id, password_reset.expires_at, users.id, users.email, users.password
        FROM password_reset
        INNER JOIN users ON users.id = password_reset.user_id
        WHERE password_reset.token_hash = $1;`, tokenHash)
	err = row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &us.ID, &us.Email, &us.Password)
	if err != nil {
		data.Token = ""

		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, ErrNotFount.Error())
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.Delete(pwReset)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if time.Now().After(pwReset.ExpiresAt) {
		return errors.Public(nil, fmt.Sprintf("%s: %s", ErrTokenExpired, data.Token))
	}

	service := user.NewService(s.ctx)
	err = service.UpdatePassword(us)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = session.NewService(s.ctx).Create(&session.Session{UserID: us.ID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) Delete(pwReset *PasswordReset) error {
	op := "model.pwreset.Delete"

	d := database.CtxDatabase(s.ctx)

	_, err := d.Exec(`
        DELETE FROM password_reset
        WHERE id = $1;`, pwReset.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
