package user

import (
	"database/sql"
	"fmt"
	"github.com/eliofery/golang-image/internal/app/models/session"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/email"
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
)

var (
	ErrEmailAlreadyExists = errors.New("email адрес уже существует")
	ErrLoginOrPassword    = errors.New("неверный логин или пароль")
)

type User struct {
	ID       uint   `validate:"omitempty"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=10,lte=32"`
}

type Service struct {
	ctx router.Ctx

	email   *email.Service
	session *session.Service
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx: ctx,

		email:   email.NewService(),
		session: session.NewService(ctx),
	}
}

func (s *Service) SignUp(mail, password string) (*User, error) {
	op := "model.user.SignUp"

	user := &User{
		Email:    strings.ToLower(mail),
		Password: password,
	}

	err := s.ctx.Validate.Struct(user)
	if err != nil {
		return user, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	row := s.ctx.DB.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email, string(hashedPassword),
	)
	err = row.Scan(&user.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return user, errors.Public(err, ErrEmailAlreadyExists.Error())
			}
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	err = s.email.Send(email.Email{
		From:    os.Getenv("EMAIL_SUPPORT"),
		To:      user.Email,
		Subject: "Регистрация на сайте",
		Plaintext: `
            Вы зарегистрировались на сайте.

            Добро пожаловать к нам на сайт.
            Приятного время провождения.

            Почта: ` + user.Email + `
            Пароль: ` + user.Password + `
        `,
		HTML: `
	       <h1>Вы зарегистрировались на сайте.</h1>

	       <p>Добро пожаловать к нам на сайт.</p>
	       <p>Приятного время провождения.

            <p><b>Почта:</b> ` + user.Email + `</p>
            <p><b>Пароль:</b> ` + user.Password + `</p>
        `,
	})
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}

	err = s.session.Create(&session.Session{UserID: user.ID})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *Service) SignIn(mail, password string) (*User, error) {
	op := "model.user.SignIn"

	user := &User{
		Email:    strings.ToLower(mail),
		Password: password,
	}

	err := s.ctx.Validate.Struct(user)
	if err != nil {
		return user, err
	}

	row := s.ctx.DB.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Public(err, ErrLoginOrPassword.Error())
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.Public(err, ErrLoginOrPassword.Error())
	}

	err = s.session.Create(&session.Session{UserID: user.ID})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *Service) UpdatePassword(us *User) error {
	op := "model.user.UpdatePassword"

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	passwordHash := string(hashedBytes)

	_, err = s.ctx.DB.Exec(`
        UPDATE users
        SET password = $1
        WHERE id = $2;`, passwordHash, us.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func GetCurrentUser(r *http.Request) (*User, error) {
	op := "model.user.CurrentUser"

	userData := &User{}

	token, err := cookie.Get(r, cookie.Session)
	if err != nil {
		return userData, fmt.Errorf("%s: %w", op, err)
	}
	tokenHash := rand.HashToken(token)

	db := database.CtxDatabase(r.Context())
	row := db.QueryRow(`
       SELECT users.id, users.email, users.password
       FROM users
       INNER JOIN sessions ON users.id = sessions.user_id
       WHERE sessions.token_hash = $1;
   `, tokenHash)
	err = row.Scan(&userData.ID, &userData.Email, &userData.Password)
	if err != nil {
		return userData, fmt.Errorf("%s: %w", op, err)
	}

	return userData, nil
}
