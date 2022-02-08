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
	DatabaseURI    = "postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable"
	MigrateDirName = "file:////home/roro/Documents/codes/borda/migrations"
	DropStmt       = "DROP SCHEMA public CASCADE;CREATE SCHEMA public;"
)

//Создает тестового юзера
func makeTestUser(db *sqlx.DB, username string, password string, contact string) (err error, user entity.User) {
	query := "INSERT INTO user (name, password, contact) VALUES($1, $2, $3)"

	err = db.QueryRowx(query, username, password, contact).Scan(&user)
	if err != nil {
		return err, user
	}
	return nil, user

}

//Создает две тестовые роли
func makeTestRole(db *sqlx.DB, roleName string) (err error, id int) {
	query := "INSERT INTO role (role_name) VALUES($1)"

	id = -1
	err = db.QueryRow(query, roleName).Scan(&id)
	if err != nil || id < 0 {
		return err, id
	}

	return nil, id
}

//Связывает роль и юзера ТЕСТОВО ЛОКАЛЬНО
func makeTestUserRoles(db *sqlx.DB, userId int, roleId int) (err error, userRoles entity.UserRoles) {
	query := "INSERT INTO user_roles (user_id, role_id) VALUES($1, $2)"

	err = db.QueryRowx(query, userId, roleId).Scan(&userRoles)
	if err != nil {
		return err, userRoles
	}

	return nil, userRoles
}

func setUp() (*sqlx.DB, interfaces.UserRepository) {
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

	var repo interfaces.UserRepository = NewPostgresUserRepository(db)

	return db, repo
}

//Done
func TestUserRepositoryCreate(t *testing.T) {
	_, repo := setUp()
	assert := assert.New(t)

	username := "roro"
	password := "12345"
	contact := "@roro"

	user, err := repo.Create(username, password, contact)

	assert.Equal(nil, err, "err should nil")
	assert.Equal(user.Username, username, "should be equal")
	assert.Equal(user.Password, password, "should be equal")
	assert.Equal(user.Contact, contact, "should be equal")

}

//Done
func TestUserRepositoryUpdatePassword(t *testing.T) {
	db, repo := setUp()
	assert := assert.New(t)

	err, user := makeTestUser(db, "roro", "1234", "@roro")
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestUser error")
	}

	password := "4321"
	err = repo.UpdatePassword(user.Id, password)

	assert.Equal(nil, err, "should be equal")
	assert.Equal(user.Password, password, "should be equal")
}

//????
func TestUserRepositoryAssignRole(t *testing.T) {
	db, repo := setUp()
	assert := assert.New(t)

	err, user := makeTestUser(db, "roro", "1234", "@roro")
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestUser error")
	}

	roleName := "master"
	err, roleId := makeTestRole(db, roleName)
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestRole error")
	}

	err = repo.AssignRole(user.Id, roleId)
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestUser error")
	}

	assert.Equal(nil, err, "should be nil")
}

func TestUserRepositoryGetRole(t *testing.T) {
	db, repo := setUp()

	roleId := 1
	roleName := "user"
	username := "roro"
	password := "12345"
	contact := "@mail"

	err, _ := makeTestRole(db, roleName)
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestRole error")
	}

	err, user := makeTestUser(db, username, password, contact)
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestUser error")
	}

	err, _ = makeTestUserRoles(db, user.Id, roleId)
	if err != nil {
		fmt.Printf("err: %v/n", err)
		panic("makeTestUserRoles error")
	}

	role, err := repo.GetRole(user.Id)

	assert := assert.New(t)
	assert.Equal(nil, err, "should be nil")
	assert.Equal(role.Id, roleId, "should be equal")
	assert.Equal(role.Name, roleName, "should be equal")

}
