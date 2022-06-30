package postgres

import (
	"borda/internal/domain"
	"database/sql"
	"errors"
	"strings"

	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) SaveUser(username, password, contact string) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}

	isUserExistQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		userTable,
	)

	var isUserExist bool
	if err := tx.Get(&isUserExist, isUserExistQuery, username); err != nil {
		return -1, err
	}

	if isUserExist {
		return -1, NewErrAlreadyExist("user", "username", username)
	}

	// Save user to the database
	createUserQuery := fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			password,
			contact
		)
		VALUES($1, $2, $3)
		RETURNING id`,
		userTable,
	)

	var userId int
	if err := tx.Get(&userId, createUserQuery, username, password, contact); err != nil {
		return -1, err
	}

	// Hardcode role to 'user'
	if err := assignRole(tx, userId, 2); err != nil {
		return -1, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return -1, err
	}

	return userId, nil
}

func (r UserRepository) GetUserByCredentials(username, password string) (*domain.User, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM public.%s
		WHERE name=$1 AND password=$2`,
		userTable,
	)

	var user domain.User
	if err := r.db.Get(&user, query, username, password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrNotFound("user", "username, password",
				strings.Join([]string{username, password}, ", "))
		}
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) GetUserById(id int) (*domain.User, error) {
	getUserQuery := fmt.Sprintf(`
		SELECT 
			u.id,
			u.name,
			u.password,
			u.contact,
			COALESCE (m.team_id, 0) AS team_id
		FROM public.%s AS u
		LEFT JOIN public.%s AS m ON u.id = m.user_id
		WHERE u.id=$1
		LIMIT 1`,
		userTable,
		teamMembersTable,
	)

	var user domain.User
	if err := r.db.Get(&user, getUserQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewErrNotFound("user", "id", id)
		}
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) GetAllUsers() ([]*domain.User, error) {
	getUsersQuery := fmt.Sprintf(`
		SELECT  
			u.id, 
			u.name, 
			u.password, 
			u.contact, 
			COALESCE (m.team_id, 0) AS team_id
		FROM public.%s AS u
		LEFT JOIN public.%s AS m ON u.id = m.user_id`,
		userTable,
		teamMembersTable)

	users := make([]*domain.User, 0)
	if err := r.db.Select(&users, getUsersQuery); err != nil {
		return nil, err
	}

	return users, nil
}

func (r UserRepository) UpdatePassword(userId int, newPassword string) error {
	query := fmt.Sprintf(`
		UPDATE public.%s
		SET password = $1
		WHERE id = $2`,
		userTable,
	)

	if _, err := r.db.Exec(query, newPassword, userId); err != nil {
		return err
	}

	return nil
}

func (r UserRepository) AssignRole(userId, roleId int) error {
	// Begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := assignRole(tx, userId, roleId); err != nil {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func assignRole(tx *sqlx.Tx, userId, roleId int) error {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			user_id,
			role_id
		)
		VALUES($1, $2)`,
		userRolesTable)

	if _, err := tx.Exec(query, userId, roleId); err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetUserRole(userId int) (*domain.Role, error) {
	query := fmt.Sprintf(`
		SELECT r.id, r.name
		FROM public.%s AS r
		INNER JOIN public.%s AS ur ON r.id=ur.role_id
		WHERE ur.user_id = $1;`,
		roleTable, userRolesTable)

	var role domain.Role
	if err := r.db.Get(&role, query, userId); err != nil {
		return nil, err
	}

	return &role, nil
}

//func (r UserRepository) IsUsernameExists(username string) error {
//	query := fmt.Sprintf(`
//		SELECT EXISTS (
//			SELECT 1
//			FROM public.%s
//			WHERE name=$1
//			LIMIT 1
//		)`,
//		r.userTableName)
//
//	var isUserExist bool
//	err := r.db.QueryRow(query, username).Scan(&isUserExist)
//	if err != nil {
//		return err
//	}
//
//	if isUserExist {
//		return domain.ErrUserAlreadyExists
//	}
//
//	return nil
//}
