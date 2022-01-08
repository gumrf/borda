package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	// Какой драйвер для базы выбрать ???
	_ "github.com/lib/pq"
	// "github.com/jackc/pgx"
	// "github.com/go-pg/pg"
)

const (
	DBPingTimeout     int = 5
	DBMaxPingAttempts int = 5
)

type DB struct {
	*sqlx.DB
	DataSourceName string
}

func NewPostgresDatabase(uri string) (*DB, error) {
	db, err := sqlx.Open("postgres", uri)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %w", err)
	}

	attempts := 0
	for {
		err := db.Ping()

		if attempts >= DBMaxPingAttempts {
			return nil, fmt.Errorf("reach max ping attempts: %w", err)
		}

		if err != nil {
			time.Sleep(time.Duration(DBPingTimeout) * time.Second)
			attempts++
			fmt.Println("Retries left", DBMaxPingAttempts-attempts, fmt.Sprintf("[Error]: %v", err))
			continue
		}

		break
	}

	return &DB{
		DB:             db,
		DataSourceName: uri,
	}, nil
}
