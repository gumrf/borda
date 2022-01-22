package repository

import (
	"borda/internal/app"
	"borda/internal/core/interfaces"
	pdb "borda/pkg/postgres"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DatabaseURI = "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable"
	MigrateDirName = "file:///home/fnc/Desktop/projects/borda/migrations/"
	DropStmt = "DROP SCHEMA public CASCADE;CREATE SCHEMA public;"
)



func setUp() (interfaces.SettingsRepository) {
	db, err := pdb.NewConnection(DatabaseURI)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Db not connected")
	}
	db.Exec(DropStmt)

	if err := app.Migrate(db, DatabaseURI, MigrateDirName); err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Failed migration")
	}

	var repo interfaces.SettingsRepository = NewPostgresSettingsRepository(db)

	return repo
}

func TestSettingsRepositorySet(t *testing.T) {
	repo := setUp()
	assert := assert.New(t)
	key := "team_limit"
	value := "5"
	settingsId := 1
	id, err := repo.Set(key, value)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	assert.Equal(id, settingsId, "they should be equal")
	assert.NotNil(id, "must be not nil")
}

func TestSettingsRepositoryGet(t *testing.T) {
	repo := setUp()
	assert := assert.New(t)
	key := "team_limit"
	value := "5"
	_, err := repo.Set(key, value)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test settings not created")
	}

	testValue, err := repo.Get(key)
	
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test settings asserted not created")
	}
	assert.Equal(value, testValue, "they should be equal")
	_, err = repo.Get("1337")
	assert.Error(err, "Settings not found with key=1337", "they should be equal")

}