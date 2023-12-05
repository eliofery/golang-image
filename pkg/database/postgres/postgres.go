package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"os"
	"strconv"
)

const (
	Port = 5432
)

type Storage struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func New() *Storage {
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		port = Port
	}

	return &Storage{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     port,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DATABASE"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

func (d *Storage) Init() (*sql.DB, error) {
	op := "postgres.Init"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
