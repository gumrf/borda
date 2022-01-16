package app

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB, databaseURI string, migrationsDirName string) (err error) {
	// driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	// if err != nil {
	// 	return err
	// }
	// m, err := migrate.NewWithDatabaseInstance(
	// 	"file:///migrations",
	// 	"postgres", driver)
	m, err := migrate.New(migrationsDirName, databaseURI)
	if err != nil {
		return fmt.Errorf("init migrations: %w", err)
	}

	// if err := m.Drop(); err != nil {
	// 	return fmt.Errorf("Drop %v", err)
	// }
	// if err := m.Force(0); err != nil {
	// 	return fmt.Errorf("Force: %w", err)
	// }
	if err := m.Up(); err != nil {
		if !errors.Is(migrate.ErrNoChange, err) {
			return fmt.Errorf("Up: %vx", err)
		}
	}

	return nil
}
