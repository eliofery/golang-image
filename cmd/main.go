package main

import (
	"database/sql"
	"github.com/eliofery/golang-image/internal/app/http/controllers/gallery"
	"github.com/eliofery/golang-image/internal/app/http/controllers/home"
	"github.com/eliofery/golang-image/internal/app/http/controllers/notfound"
	"github.com/eliofery/golang-image/internal/app/http/controllers/user"
	mw "github.com/eliofery/golang-image/internal/app/http/middleware"
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/database/postgres"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
)

func main() {
	// Подключение конфигурационного файла .env
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Подключение логирования
	logger := logging.New()

	// Подключение к БД Postgres
	db, err := database.Init(postgres.New())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(db)

	// Миграция БД Postgres
	if err = database.MigrateFS(db, goose.DialectPostgres); err != nil {
		logger.Error(err.Error())
	}

	// Подключение валидатора
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Создание роутера
	route := router.New()

	// Пользовательский Middleware
	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mw.Csrf)
	route.Use(mw.Inject(logger, db, validate))
	route.Use(mw.SetUser)

	// Подключение ресурсов
	tpl.AssetsFsInit(route.Chi)

	// Роуты
	route.Get("/", home.Index)

	route.NotFound(notfound.Page404)
	route.MethodNotAllowed(notfound.Page405)

	route.Group(func(r *router.Router) {
		r.Chi.Use(mw.Auth)

		r.Route("/user", func(r *router.Router) {
			r.Get("/", user.Index)
			r.Post("/logout", user.Logout)
		})

		r.Route("/gallery", func(r *router.Router) {
			r.Get("/new", gallery.New)
		})
	})

	route.Group(func(r *router.Router) {
		r.Chi.Use(mw.Guest)

		r.Get("/signup", user.SignUp)
		r.Post("/signup", user.Create)

		r.Get("/signin", user.SignIn)
		r.Post("/signin", user.Auth)

		r.Get("/forgot-pw", user.ForgotPassword)
		r.Post("/forgot-pw", user.ProcessForgotPassword)

		r.Get("/reset-pw", user.ResetPassword)
		r.Post("/reset-pw", user.ProcessResetPassword)
	})

	// Запуск сервера
	logger.Info("Сервер запущен: http://localhost:8080")
	err = http.ListenAndServe(":8080", route.ServeHTTP())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
