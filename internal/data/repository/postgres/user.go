package postgres

import (
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

var _ interfaces.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

//Создать нового юзера в бд
func (user *PostgresUserRepository) Create(username, passwordHash, contact string, roleId int) (int, error) {
	qwery := "INSERT INTO user (name, password, contact) VALUES(?, ?, ?)"

	result, err := user.db.Exec(qwery, username, passwordHash, contact)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (user *PostgresUserRepository) UpdatePassword(newHashPassword string, id int) error {
	qwery := "UPDATE user SET password = $1 WHERE id = $2"

	_, err := user.db.Exec(qwery, newHashPassword, id)
	if err != nil {
		return err
	}
	return nil
}
