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

func (r PostgresUserRepository) Create(username, password, contact string) (user entity.User, err error) {
	query := "INSERT INTO user (name, password, contact) VALUES(?, ?, ?)"

	err = r.db.QueryRow(query, username, password, contact).Scan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r PostgresUserRepository) UpdatePassword(userId int, newPassword string) error {
	query := "UPDATE user SET password = $1 WHERE id = $2"

	_, err := r.db.Exec(query, newPassword, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r PostgresUserRepository) AssignRole(userId, roleId int) error {
	query := "INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)"

	_, err := r.db.Exec(query, userId, roleId)
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresUserRepository) GetRole(userId int) (role entity.Role, err error) {
	query := "SELECT * FROM user_roles INNER JOIN role ON user_roles.role_id = role.role_id WHERE user_roles.user_id = $1"

	err = r.db.QueryRowx(query, userId).Scan(&role)
	if err != nil {
		return role, err
	}
	return role, nil
}
