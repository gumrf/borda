package postgres

import (
	"errors"
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

var (
	ErrDatabasePingTimeout = errors.New("reach ping timeout")
)

type DB struct {
	*sqlx.DB
}

// Options contains postgres database connection options.
// type Options struct {
// 	Host string
// 	Port int

// 	User     string
// 	Password string
// 	Database string

// 	SSLMode string
// }

func NewPostgresDatabase(uri, username, password string) (*DB, error) {

	// opts, err := pq.ParseURL(uri)
	// db := pg.Connect(opts)

	// opts, err := pq.NewConnector(uri)
	// if err != nil {

	// }

	// dataSourceName := pq.ParseURL(url)
	// connector := pq.NewConnector(uri)
	// if username != "" && password != "" {
	// 	opts.SetAuth(options.Credential{
	// 		Username: username, Password: password,
	// 	})
	// }

	// dataSourceName := fmt.Sprintf(
	// 	"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
	// 	opts.Host, opts.Port, opts.User, opts.Password, opts.Database, opts.SSLMode,
	// )

	// TODO: if specified username and password parse uri and override values

	db, err := sqlx.Open("postgres", uri)
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

		break
	}

	return &DB{db}, nil
}
