package repository

import (
	"github.com/jmoiron/sqlx"
)

type PostrgresRoleRepository struct {
	db *sqlx.DB
}

// var _ interfaces.RoleRepository = (*PostrgresRoleRepository)(nil)

func newPostrgresRoleRepository(db *sqlx.DB) *PostrgresRoleRepository {
	return &PostrgresRoleRepository{db: db}
}

func (r PostrgresRoleRepository) Get(id int) (string, error) {
	qwery := "SELECT role_id FROM user_roles WHERE user_Id = $1"

	var a int
	err := r.db.QueryRowx(qwery, id).Scan(&a)
	if err != nil {
		return "", err
	}

	qwery = "SELECT role_name FROM role WHERE Id = $1"

	var b string
	err = r.db.QueryRowx(qwery, a).Scan(&b)
	if err != nil {
		return "", err
	}

	return b, nil
}

func (r PostrgresRoleRepository) Give(userId, roleId int) error {
	qwery := "INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)"

	_, err := r.db.Exec(qwery, userId, roleId)
	if err != nil {
		return err
	}
	return nil
}
