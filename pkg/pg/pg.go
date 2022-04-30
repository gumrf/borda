package pg

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// timeout is Postgres timeout when trying to connect
	timeout = 5
)

func Open(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	retries := 5
	for {
		// Try to ping Postgres DB
		if err := db.Ping(); err != nil {
			if retries >= 0 {
				fmt.Printf("%v [retries left: %d]\n", err, retries)
				retries--
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			return nil, errors.New("connecting to Postgres failed after maximum attempts")

		}
		break
	}
	// db.SetMaxIdleConns(idleConn)
	// db.SetMaxOpenConns(maxConn)
	return db, nil
}

func Migrate(db *sqlx.DB, migrationsPath string, version uint) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("connect db driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialize migrations: %w", err)
	}

	if err := m.Migrate(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

func Close(db *sqlx.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("close connection to Postgres: %w", err)
	}

	return nil
}
