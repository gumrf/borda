package postgres

import (
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostrgresRoleRepository struct {
	db *sqlx.DB
}

var _ interfaces.RoleRepository = (*PostrgresRoleRepository)(nil)

func NewPostrgresRoleRepository(db *sqlx.DB) *PostrgresRoleRepository {
	return &PostrgresRoleRepository{db: db}
}

func (role *PostrgresRoleRepository) Get(id int) (int, error) {
	qwery := "SELECT role_id FROM user_roles WHERE user_Id = $1"

	var a int
	err := role.db.QueryRowx(qwery, id).Scan(&a)
	if err != nil {
		return -1, err
	}
	return a, nil
}

func (role *PostrgresRoleRepository) GiveRole(userId, roleId int) error {
	qwery := "INSERT INTO user_roles (user_id, role_id) VALUES(?, ?)"

	_, err := role.db.Exec(qwery, userId, roleId)
	if err != nil {
		return err
	}
	return nil
}
