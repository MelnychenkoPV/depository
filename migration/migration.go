package migration

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	pgx5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func NewMigration(db *sql.DB) (*migrate.Migrate, error) {
	srcDrv, err := iofs.New(fs, ".")
	if err != nil {
		return nil, err
	}

	dbDrv, err := pgx5.WithInstance(db, new(pgx5.Config))
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("iofs", srcDrv, "pgx5", dbDrv)
}
