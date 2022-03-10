package postgres

import (
	"borda/internal/core"
	"borda/internal/core/entity"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db                 *sqlx.DB
	userTableName      string
	roleTableName      string
	userRolesTableName string
}

var _ core.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sqlx.DB) core.UserRepository {
	return UserRepository{
		db:                 db,
		userTableName:      "\"user\"",
		roleTableName:      "role",
		userRolesTableName: "user_role",
	}
}

func (r UserRepository) Create(username, password, contact string) (userId int, err error) {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			password,
			contact
		) 
		VALUES($1, $2, $3)
		RETURNING id`,
		r.userTableName)

	var id int = -1
	err = r.db.QueryRowx(query, username, password, contact).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (r UserRepository) UpdatePassword(userId int, newPassword string) error {
	query := fmt.Sprintf(`
		UPDATE public.%s
		SET password = $1
		WHERE id = $2`,
		r.userTableName)

	_, err := r.db.Exec(query, newPassword, userId)
	if err != nil {
		return fmt.Errorf("Can't update password: %w", err)
	}

	return nil
}

func (r UserRepository) AssignRole(userId, roleId int) error {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			user_id,
			role_id
		)
		VALUES($1, $2)`,
		r.userRolesTableName)

	_, err := r.db.Exec(query, userId, roleId)
	if err != nil {
		return fmt.Errorf("Can't add role: %w", err)
	}

	return nil
}

func (r UserRepository) GetRole(userId int) (entity.Role, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.name
		FROM public.%s AS r
		INNER JOIN public.%s AS ur ON r.id=ur.role_id
		WHERE ur.user_id = $1;`,
		r.roleTableName, r.userRolesTableName)

	var role entity.Role
	err := r.db.Get(&role, query, userId)
	if err != nil {
		return entity.Role{}, err
	}

	return role, nil
}
