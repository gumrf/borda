package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

var _ interfaces.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(db *sqlx.DB) interfaces.UserRepository {
	return PostgresUserRepository{db: db}
}

func (r PostgresUserRepository) Create(username, password, contact string) (userId int, err error) {
	qwery := "INSERT INTO user (name, password, contact) VALUES(?, ?, ?)"

	result, err := r.db.Exec(qwery, username, password, contact)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (r PostgresUserRepository) UpdatePassword(userId int, newPassword string) error {
	qwery := "UPDATE user SET password = $1 WHERE id = $2"

	_, err := r.db.Exec(qwery, newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresUserRepository) AssignRole(userId, roleId int) error {
	qwery := "INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)"

	_, err := r.db.Exec(qwery, userId, roleId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresUserRepository) GetRole(userId int) (role entity.Role, err error) {
	qwery := "SELECT role_id FROM user_roles WHERE user_id = $1"

	var id int
	err = r.db.QueryRowx(qwery, userId).Scan(&id)
	if err != nil {
		return role, err
	}

	qwery = "SELECT role_name FROM role WHERE Id = $1"

	err = r.db.QueryRowx(qwery, id).Scan(&role)
	if err != nil {
		return role, err
	}

	return role, nil
}
