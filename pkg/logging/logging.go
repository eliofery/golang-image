package logging

import (
	"log/slog"
	"os"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

func New() *slog.Logger {
	var slHandler slog.Handler

	switch os.Getenv("ENV") {
	case Dev:
		slHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case Prod:
		slHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(slHandler)
}
