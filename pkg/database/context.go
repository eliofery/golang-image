package database

import (
	"database/sql"
	"golang.org/x/net/context"
)

type key string

const databaseKey key = "database"

func WithDatabase(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, databaseKey, db)
}

func CtxDatabase(ctx context.Context) *sql.DB {
	val := ctx.Value(databaseKey)

	db, ok := val.(*sql.DB)
	if !ok {
		return nil
	}

	return db
}
