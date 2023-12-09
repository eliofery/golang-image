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
	"github.com/go-playground/validator/v10"
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
	ctx      router.Ctx
	db       *sql.DB
	validate *validator.Validate
	email    *email.Service

	session *session.Service
}

func NewService(ctx router.Ctx) *Service {
	return &Service{
		ctx:      ctx,
		db:       ctx.DB,
		validate: ctx.Validate,
		email:    email.NewService(),

		session: session.NewService(ctx),
	}
}

func (s *Service) SignUp(us *User) error {
	op := "model.us.SignUp"

	err := s.validate.Struct(us)
	if err != nil {
		return err
	}

	us.Email = strings.ToLower(us.Email)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		us.Email, string(hashedPassword),
	)
	err = row.Scan(&us.ID)
	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) {
			if pgError.Code == pgerrcode.UniqueViolation {
				return errors.Public(err, ErrEmailAlreadyExists.Error())
			}
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.email.Send(email.Email{
		From:    os.Getenv("EMAIL_SUPPORT"),
		To:      us.Email,
		Subject: "Регистрация на сайте",
		Plaintext: `
            Вы зарегистрировались на сайте.

            Добро пожаловать к нам на сайт.
            Приятного время провождения.

            Почта: ` + us.Email + `
            Пароль: ` + us.Password + `
        `,
		HTML: `
	       <h1>Вы зарегистрировались на сайте.</h1>

	       <p>Добро пожаловать к нам на сайт.</p>
	       <p>Приятного время провождения.

            <p><b>Почта:</b> ` + us.Email + `</p>
            <p><b>Пароль:</b> ` + us.Password + `</p>
        `,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.session.Create(&session.Session{UserID: us.ID})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(user *User) error {
	op := "model.user.SignIn"

	err := s.validate.Struct(user)
	if err != nil {
		return err
	}

	user.Email = strings.ToLower(user.Email)
	password := user.Password

	row := s.db.QueryRow("SELECT * FROM users WHERE email = $1", user.Email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Public(err, ErrLoginOrPassword.Error())
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.Public(err, ErrLoginOrPassword.Error())
	}

	err = s.session.Create(&session.Session{UserID: user.ID})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdatePassword(us *User) error {
	op := "model.user.UpdatePassword"

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	passwordHash := string(hashedBytes)

	_, err = s.db.Exec(`
        UPDATE users
        SET password = $2
        WHERE id = $1;`, us.ID, passwordHash)
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
