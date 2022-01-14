package postgres

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
	"github.com/jmoiron/sqlx"
)

func NewConnection(uri string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", uri)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	retries := 5
	for {
		// Try to ping database.
		if err := db.Ping(); err != nil {
			defer db.Close() // close database connection

			if retries > 0 {
				fmt.Println("Retries left", retries, fmt.Sprintf("[Error]: %v", err))
				retries--
				time.Sleep(time.Duration(5) * time.Second)
				continue
			}

			return nil, fmt.Errorf("can't ping database %w", err)
		}

		break
	}

	return db, nil
}
