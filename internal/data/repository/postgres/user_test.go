package postgres_test

import (
	"borda/internal/core/entity"
	"borda/internal/data/repository/postgres"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func makeTestRole(t *testing.T, db *sqlx.DB, roleName string) int {
	query := "INSERT INTO role (name) VALUES($1) RETURNING id"

	id := -1
	err := db.Get(&id, query, roleName)
	if err != nil {
		t.Fatalf("makeTestRole error: %v", err)
	}

	if id < 1 {
		t.Fatal("makeTestRole error: id must be > 0")
	}

	return id
}

func Test_UserRepository_Create(t *testing.T) {
	db := MustConnectAndMigrate(t)
	repo := postgres.NewUserRepository(db)

	user := entity.User{
		Username: "roro",
		Password: "12345",
		Contact:  "@roro",
	}
	username := "roro"
	password := "12345"
	contact := "@roro"

	userId, err := repo.Create(username, password, contact)

	assert := assert.New(t)
	assert.Equal(nil, err, "err should nil")
	assert.Equal(1, userId)

	var testedUser entity.User
	if err := db.Get(&testedUser, `SELECT * FROM public."user" WHERE id=$1`, userId); err != nil {
		t.Fatalf("get user error: %v", err)
	}

	assert.Equal(testedUser.Username, user.Username, "should be equal")
	assert.Equal(testedUser.Password, user.Password, "should be equal")
	assert.Equal(testedUser.Contact, user.Contact, "should be equal")
}

func Test_UserRepository_UpdatePassword(t *testing.T) {
	db := MustConnectAndMigrate(t)
	repo := postgres.NewUserRepository(db)
	assert := assert.New(t)

	userId, err := makeTestUser(db, "roro", "1234", "@roro")
	if err != nil {
		t.Fatalf("makeTestUser error%v\n", err)
	}

	newPassword := "4321"
	err = repo.UpdatePassword(userId, newPassword)

	assert.Equal(nil, err, "should be nil")

	var password string
	if err := db.Get(&password, `SELECT password FROM public."user" WHERE id=$1`, userId); err != nil {
		t.Fatalf("get user password error: %v", err)
	}

	assert.Equal(newPassword, password, "should be equal")
}

func Test_UserRepository_AssignRole(t *testing.T) {
	db := MustConnectAndMigrate(t)
	repo := postgres.NewUserRepository(db)
	assert := assert.New(t)

	userId, err := makeTestUser(db, "roro", "1234", "@roro")
	if err != nil {
		t.Fatalf("makeTestUser error: %v/n", err)
	}

	roleName := "master"
	roleId := makeTestRole(t, db, roleName)

	err = repo.AssignRole(userId, roleId)

	assert.Equal(nil, err, "Error should be nil")
}

func Test_UserRepository_GetRole(t *testing.T) {
	db := MustConnectAndMigrate(t)
	repo := postgres.NewUserRepository(db)

	user := entity.User{
		Id:       1,
		Username: "roro",
		Password: "12345",
		Contact:  "@mail",
	}

	roleName := "admin"
	roleId := makeTestRole(t, db, roleName)

	userId, err := makeTestUser(db, user.Username, user.Password, user.Contact)
	if err != nil {
		t.Fatalf("makeTestUser error: %v/n", err)
	}

	// make sure method work properly or check error
	repo.AssignRole(userId, roleId) // nolint

	role, err := repo.GetRole(user.Id)

	assert := assert.New(t)
	assert.Equal(nil, err, "should be nil")
	assert.Equal(role.Id, roleId, "should be equal")
	assert.Equal(role.Name, roleName, "should be equal")
}
