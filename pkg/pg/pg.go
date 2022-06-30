package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	timeout = 5
)

type Options struct {
	Logger log.Logger
}

func Connect(url string) (*pgxpool.Pool, error) {
	pgxpoolConfig, _ := pgxpool.ParseConfig(url)
	pgxpoolConfig.ConnConfig.PreferSimpleProtocol = true
	// pgxpoolConfig.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), pgxpoolConfig)
	if err != nil {
		return nil, fmt.Errorf("can't create pgx pool: %w", err)
	}

	retries := 5
	for {
		// Try to ping Postgres DB
		if err := pool.Ping(context.Background()); err != nil {
			if _, ok := err.(net.Error); !ok {
				return nil, err
			}

			if retries >= 0 {
				fmt.Printf("%v [retries left: %d]\n", err, retries)
				retries--
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			return nil, errors.New("reach the maximum number of attempts")
		}

		break
	}

	return pool, nil
}

func Migrate(connectionURL, source string, version uint) error {
	m, err := migrate.New(
		source, strings.Replace(connectionURL, "postgres", "pgx", 1),
	)
	if err != nil {
		return err
	}

	// m.Drop()

	migrateError := m.Migrate(version)
	if migrateError != nil && migrateError != migrate.ErrNoChange {
		return fmt.Errorf("can't apply migrations: %v", migrateError)
	}

	return nil
}
