package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	DBPingTimeout     int = 5
	DBPingMaxAttempts int = 5
)

type DB struct {
	*sqlx.DB
}

// Options contains postgres database connection options.
type Options struct {
	Host string
	Port int

	User     string
	Password string
	Database string

	SSLMode string
}

func NewPostgresDatabase(opts Options) (*DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		opts.Host, opts.Port, opts.User, opts.Password, opts.Database, opts.SSLMode,
	)

	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Open: %v", err)
	}

	retries := 0
	for {
		if retries > DBPingMaxAttempts {
			return nil, fmt.Errorf("reach retries limit")
		}

		err := db.Ping()
		if err != nil {
			log.Println(err, "RETRIES LEFT:", DBPingMaxAttempts-retries)
			time.Sleep(time.Duration(DBPingTimeout) * time.Second)
			retries++
			continue
		}

		log.Println("Connected to database")
		break
	}

	return &DB{db}, nil
}
