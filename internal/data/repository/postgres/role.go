package postgres

import (
	"borda/internal/core/interfaces"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type PostrgresRoleRepository struct {
	db *sqlx.DB
}

var _ interfaces.RoleRepository = (*PostrgresRoleRepository)(nil)

func NewPostrgresRoleRepository(db *sqlx.DB) *PostrgresRoleRepository {
	return &PostrgresRoleRepository{db: db}
}

func (role *PostrgresRoleRepository) Get(id int) (sql.Result, error) {
	qwery := "SELECT RoleId FROM User_roles WHERE UserId = $1"

	result, err := role.db.Exec(qwery, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (role *PostrgresRoleRepository) Associate(userId, roleId int) error {
	qwery := "INSERT INTO User_roles (UserId, RoleId) VALUES(?, ?)"

	_, err := role.db.Exec(qwery, userId, roleId)
	if err != nil {
		return err
	}
	return nil
}
