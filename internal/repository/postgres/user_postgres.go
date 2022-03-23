package postgres

import (
	"borda/internal/domain"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db                 *sqlx.DB
	userTableName      string
	roleTableName      string
	userRolesTableName string
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:                 db,
		userTableName:      "\"user\"",
		roleTableName:      "role",
		userRolesTableName: "user_role",
	}
}

// TODO: pass user object when create user
func (r UserRepository) CreateNewUser(username, password, contact string) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}

	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		r.userTableName)

	var isUserExist bool
	err = tx.QueryRow(query, username).Scan(&isUserExist)
	if err != nil {
		return -1, err
	}

	if isUserExist {
		return -1, domain.ErrUserAlreadyExists
	}

	query = fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			password,
			contact
		)
		VALUES($1, $2, $3)
		RETURNING id`,
		r.userTableName,
	)

	var userId int
	err = tx.Get(&userId, query, username, password, contact)

	if err != nil {
		return -1, err
	}

	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return userId, nil
}

func (r UserRepository) IsUsernameExists(username string) error {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		r.userTableName)

	var isUserExist bool
	err := r.db.QueryRow(query, username).Scan(&isUserExist)
	if err != nil {
		return err
	}

	if isUserExist {
		return domain.ErrUserAlreadyExists
	}

	return nil
}

func (r UserRepository) FindUserByCredentials(username, password string) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM public.%s
		WHERE name=$1 AND password=$2`,
		r.userTableName)

	var user domain.User
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
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

func (r UserRepository) GetRole(userId int) (domain.Role, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.name
		FROM public.%s AS r
		INNER JOIN public.%s AS ur ON r.id=ur.role_id
		WHERE ur.user_id = $1;`,
		r.roleTableName, r.userRolesTableName)

	var role domain.Role
	err := r.db.Get(&role, query, userId)
	if err != nil {
		return domain.Role{}, err
	}

	return role, nil
}
