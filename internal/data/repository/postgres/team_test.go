package repository

import (
	"borda/internal/core/interfaces"
	"borda/internal/app"
	pdb "borda/pkg/postgres"
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

const (
	DatabaseURI = "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable"
)

func TestTeamRepository(t *testing.T) {
	assert := assert.New(t)
	_ = assert

	db, err := pdb.NewConnection(DatabaseURI)
	if err != nil {
		fmt.Printf("Db not connected")
	}
	if err := app.Migrate(db); err != nil {
		fmt.Printf("Failed migration: %v", err)
	}
	fmt.Printf("Migration did run successfully")

	var repo interfaces.TeamRepository = NewPostgresTeamRepository(db)
	repo.Create(1, "ShrekTeam")

	//   // assert equality
	//   assert.Equal(123, 123, "they should be equal")

	//   // assert inequality
	//   assert.NotEqual(123, 456, "they should not be equal")

	//   // assert for nil (good for errors)
	//   assert.Nil(object)

	//   // assert for not nil (good when you expect something)
	//   if assert.NotNil(object) {

	//     // now we know that object isn't nil, we are safe to make
	//     // further assertions without causing any errors
	//     assert.Equal("Something", object.Value)
	//   }
}
