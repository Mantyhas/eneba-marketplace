package database

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func Open(ctx context.Context, url string) (*sql.DB, error) {
	db, e := sql.Open("pgx", url)
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	if e = db.PingContext(ctx); e != nil {
		db.Close()
		return nil, e
	}
	return db, nil
}
