package middleware

import (
	"database/sql"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func Inject(logger *slog.Logger, db *sql.DB, validator *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.WithLogging(r.Context(), logger)
			ctx = database.WithDatabase(ctx, db)
			ctx = validate.WithValidation(ctx, validator)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
