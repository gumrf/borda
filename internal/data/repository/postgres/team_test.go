package repository

import (
	"borda/internal/app"
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"
	pdb "borda/pkg/postgres"
	"fmt"
	"testing"

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

func uploadTeamSettings(db *sqlx.DB) error {
	query := `
	INSERT INTO manage_settings
	(key, value)
	VALUES($1, $2)
	`
	_, err := db.Query(query, "team_limit", "5")
	if err != nil {
		return err
	}
	return nil
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

	if err := uploadTeamSettings(db); err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Failed load team settings")
	}

	var repo interfaces.TeamRepository = NewPostgresTeamRepository(db)

	return db, repo
}

func TestTeamRepositoryCreate(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	userId, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}

	// tests

	// Default
	teamName := "ShrekTeam"
	team, err := repo.Create(userId, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	assert.Equal(team.Name, teamName, "they should be equal")
	assert.Equal(team.TeamLeaderId, userId, "they should be equal")
	assert.NotNil(team.Token, "must be not nil")
	assert.NotNil(team.Id, "must be not nil")

	// Duplicate name
	team, err = repo.Create(userId, teamName)
	assert.Error(err)
}

func TestTeamRepositoryGet(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	userId, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}

	teamName := "ShrekTeam"
	team, err := repo.Create(userId, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team not created")
	}
	
	// tests

	// Default
	teamAssert, err := repo.Get(team.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team asserted not created")
	}
	assert.Equal(team.Name, teamAssert.Name, "they should be equal")
	assert.Equal(team.TeamLeaderId, teamAssert.TeamLeaderId, "they should be equal")
	assert.Equal(team.Token, teamAssert.Token, "they should be equal")
	assert.Equal(team.Id, teamAssert.Id, "they should be equal")

	// Not found
	teamAssert, err = repo.Get(1337)
	assert.Error(err, "Team not found with id=1337", "they should be equal")
}

func TestTeamRepositoryAddMember(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	userId, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}
	teamName := "ShrekTeam"
	team, err := repo.Create(userId, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team not created")
	}

	// tests

	// Default
	err = repo.AddMember(team.Id, userId)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}


	// Duplicate
	err = repo.AddMember(team.Id, userId)
	assert.Error(err, "User id=1 already in team with id=1", "they should be equal")
	
	// Not user
	err = repo.AddMember(team.Id, 1337)
	assert.Error(err, "User with id=%v not found", "they should be equal")

	// Not team
	err = repo.AddMember(1337, 1)
	assert.Error(err, "Team with id=%v not found", "they should be equal")
}

func TestTeamRepositoryGetMembers(t *testing.T) {
	// setup
	db, repo := setUp()
	assert := assert.New(t)
	userId, err := makeTestUser(db, "jayse", "test", "@jaysess")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}
	userId2, err := makeTestUser(db, "jayseClone", "test", "@jaysessClone")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test user not created")
	}
	teamName := "ShrekTeam"
	team, err := repo.Create(userId, teamName)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		panic("Test team not created")
	}

	err = repo.AddMember(team.Id, userId)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	err = repo.AddMember(team.Id, userId2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	// test

	// Default
	users, err := repo.GetMembers(team.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	
	assert.True(len(users) == 2)
	var _user entity.User = users[0]
	assert.Equal(_user.Username, "jayse")
	assert.Equal(_user.Password, "test")
	assert.Equal(_user.Contact, "@jaysess")
	assert.Equal(_user.Id, 1)
	assert.Equal(_user.TeamId, 1)

	_user = users[1]
	assert.Equal(_user.Username, "jayseClone")
	assert.Equal(_user.Password, "test")
	assert.Equal(_user.Contact, "@jaysessClone")
	assert.Equal(_user.Id, 2)
	assert.Equal(_user.TeamId, 1)

	// Not found team by id
	_, err = repo.GetMembers(1337)
	assert.Error(err, "Team with id=%v not found", "they should be equal")
}
