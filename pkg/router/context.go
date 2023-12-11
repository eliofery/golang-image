package router

import (
	"context"
	"database/sql"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/validate"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type Ctx struct {
	context.Context
	http.ResponseWriter
	*http.Request
	*slog.Logger
	*validator.Validate
	*sql.DB
}

func CtxRouter(ctx context.Context) Ctx {
	return Ctx{
		Context:        ctx,
		ResponseWriter: ResponseWriter(ctx),
		Request:        Request(ctx),
		Logger:         logging.Logging(ctx),
		Validate:       validate.Validation(ctx),
		DB:             database.CtxDatabase(ctx),
	}
}
