package logging

import (
	"golang.org/x/net/context"
	"log/slog"
)

type key string

const loggingKey key = "logging"

func WithLogging(ctx context.Context, logging *slog.Logger) context.Context {
	return context.WithValue(ctx, loggingKey, logging)
}

func Logging(ctx context.Context) *slog.Logger {
	val := ctx.Value(loggingKey)

	logging, ok := val.(*slog.Logger)
	if !ok {
		return nil
	}

	return logging
}
