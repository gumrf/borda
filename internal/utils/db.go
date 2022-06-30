package utils

import (
	"borda/pkg/pg"
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jmoiron/sqlx"
)

const (
	dsn            string = "postgres://postgres:secret_password@127.0.0.1:5432/borda?sslmode=disable"
	migrationsPath string = "file://../../migrations"
	dropStmt       string = "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
)

func MustOpenDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	db, err := pg.Connect(dsn)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(context.Background(), dropStmt); err != nil {
		t.Fatal(err)
	}

	if err := pg.Migrate(dsn, migrationsPath, 2); err != nil {
		t.Fatal(err)
	}

	return db
}

// MustCloseDB closes the DB. Fatal on error.
func MustCloseDB(t *testing.T, db *sqlx.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
