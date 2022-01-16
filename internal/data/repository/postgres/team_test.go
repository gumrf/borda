package repository

import (
	"borda/internal/core/interfaces"
	"borda/internal/app"
	pdb "borda/pkg/postgres"
	"testing"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const (
	DatabaseURI = "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable"
	MigrateDirName = "file:///home/jayse/Desktop/projects/borda/migrations/"
	DropStmt = "DROP SCHEMA public CASCADE;CREATE SCHEMA public;"
)

func makeTestUser(db *sqlx.DB, username string, password string, contact string) (int, error) {
	query := `
	INSERT INTO public."user"
	(name, password, contact)
	VALUES($1, $2, $3)
	RETURNING id
	`
	id := -1
	err := db.QueryRow(query, username, password, contact).Scan(&id)
	if err != nil {
		return -1, err
	}
	if id == -1 {
		return -1, err
	}
	return int(id), err
}

func setUp() (*sqlx.DB, interfaces.TeamRepository) {
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
	var repo interfaces.TeamRepository = NewPostgresTeamRepository(db)

	return db, repo
}

func TestTeamRepositoryCreate(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	user_id, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}

	// test
	teamName := "ShrekTeam"
	team, err := repo.Create(user_id, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	assert.Equal(team.Name, teamName, "they should be equal")
	assert.Equal(team.TeamLeaderId, user_id, "they should be equal")
	assert.NotNil(team.Token, "must be not nil")
	assert.NotNil(team.Id, "must be not nil")
}

func TestTeamRepositoryGet(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	user_id, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}

	teamName := "ShrekTeam"
	team, err := repo.Create(user_id, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team not created")
	}
	
	// test
	teamAssert, err := repo.Get(team.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team asserted not created")
	}
	assert.Equal(team.Name, teamAssert.Name, "they should be equal")
	assert.Equal(team.TeamLeaderId, teamAssert.TeamLeaderId, "they should be equal")
	assert.Equal(team.Token, teamAssert.Token, "they should be equal")
	assert.Equal(team.Id, teamAssert.Id, "they should be equal")
}

func TestTeamRepositoryAddMember(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	user_id, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}
	teamName := "ShrekTeam"
	team, err := repo.Create(user_id, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team not created")
	}

	// tests

	// Default
	err = repo.AddMember(team.Id, user_id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}


	// Duplicate
	err = repo.AddMember(team.Id, user_id)
	assert.Error(err, "User id=1 already in team with id=1", "they should be equal")
	
	// Not user
	err = repo.AddMember(team.Id, 1337)
	assert.Error(err, "User with id=%v not found", "they should be equal")

	// Not team
	err = repo.AddMember(1337, 1)
	assert.Error(err, "Team with id=%v not found", "they should be equal")

}
