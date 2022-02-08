package postgres_test

import (
	"borda/internal/data/repository/postgres"
	"log"
	"testing"

	"github.com/jmoiron/sqlx"
)

const (
	dsn            string = "postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable"
	migrationsPath string = "file://migrations"
	dropStmt              = "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
)

// Ensure the Postgres database can open & close.
func Test_Connect(t *testing.T) {
	db, err := postgres.Connect(dsn)
	if err != nil {
		log.Fatal("Connecting to DB failed: ", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal("Closing DB connection failed: ", err)
	}
}

func MustConnectAndMigrate(t *testing.T) *sqlx.DB {
	db, err := postgres.Connect(dsn)
	if err != nil {
		t.Fatal("Oops.. can't connect to DB: ", err)
	}

	if _, err := db.Exec(dropStmt); err != nil {
		t.Fatal("Oops.. reset shema failed: ", err)
	}

	err = postgres.RunMigrations(db, migrationsPath)
	if err != nil {
		t.Fatal("Oops.. migrations failed: ", err)
	}

	return db
}
