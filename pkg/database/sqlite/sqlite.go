package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	Path string
}

func New() *Storage {
	return &Storage{
		Path: os.Getenv("SQLITE_PATH"),
	}
}

func (d *Storage) Init() (*sql.DB, error) {
	op := "sqlite.Init"

	db, err := sql.Open("sqlite3", d.Path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
