package repository

import (
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

var _ interfaces.UserRepository = (*PostgresUserRepository)(nil)

func newPostgresUserRepository(db *sqlx.DB) PostgresUserRepository {
	return PostgresUserRepository{db: db}
}

// Create создает нового юзера в базе данных
func (r PostgresUserRepository) Create(username, password, contact string) (int, error) {
	qwery := "INSERT INTO user (name, password, contact) VALUES(?, ?, ?)"

	result, err := r.db.Exec(qwery, username, password, contact)
	if err != nil {
		return -1, err
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(userId), nil
}

func (r PostgresUserRepository) UpdatePassword(userId int, newPassword string) error {
	qwery := "UPDATE user SET password = $1 WHERE id = $2"

	_, err := r.db.Exec(qwery, newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresUserRepository) RequestRole(userId, roleId int) error
func (r PostgresUserRepository) GetRole(userId int) (roleId int, err error)
