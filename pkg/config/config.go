package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	ErrNotFound = errors.New("не удалось загрузить файл .env")
)

const (
	Dev = ".env"
)

func Init() error {
	op := "config.Init"

	var env string

	if len(os.Args) == 1 {
		env = Dev
	} else {
		env = os.Args[1]
	}

	if err := godotenv.Load(env); err != nil {
		return fmt.Errorf("%s: %w", op, ErrNotFound)
	}

	return nil
}
