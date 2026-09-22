package cmd

import (
	"context"
	"database/sql"
)

func CreateDB(ctx context.Context, cnf DB) (*sql.DB, error) {
	db, err := sql.Open("pgx", cnf.DSN)
	if err != nil {
		return nil, err
	}
	return db, db.PingContext(ctx)
}
