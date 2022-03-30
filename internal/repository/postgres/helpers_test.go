package postgres_test

import (
	"borda/pkg/pg"

	"testing"

	"github.com/jmoiron/sqlx"
)

const (
	dsn            string = "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
	migrationsPath string = "file://../../../migrations"
	dropStmt       string = "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
)

func MustOpenDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := pg.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.Exec(dropStmt); err != nil {
		t.Fatal(err)
	}

	if err := pg.Migrate(db, migrationsPath); err != nil {
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
