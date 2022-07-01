package user

import (
	shema "borda/internal/pkg"
	"borda/internal/utils"

	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
)

// const (
// 	errUserAlreadyExists = "user with name %v already exists"
// 	errUserNotFound      = "user with id %v not found"
// )

var (
	ErrUserAlreadyExists = errors.New("user already exist")
	ErrUserNotFound      = errors.New("user not found")
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r PostgresUserRepository) Save(ctx context.Context, user User) (User, error) {
	var err error

	query := fmt.Sprintf(`
INSERT INTO public.%s (
	name,
	password,
	contact
)
VALUES($1, $2, $3)
RETURNING id`, shema.Tables.User,
	)

	tx, ok := utils.ExtractTransaction(ctx)

	switch ok {
	case true:
		err = tx.GetContext(
			ctx, &user.Id, query,
			user.Username, user.Password, user.Contact,
		)
	default:
		err = r.db.GetContext(
			ctx, &user.Id, query,
			user.Username, user.Password, user.Contact,
		)
	}

	if err != nil {
		if pqerr, ok := err.(*pq.Error); ok {
			switch pqerr.Code.Name() {
			case "unique_violation":
				return User{}, fmt.Errorf(
					"%w [name=%v]", ErrUserAlreadyExists, user.Username,
				)
			}
		}
		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) SaveAll(entities []User) ([]User, error) {
	saveError := utils.Transactional(
		context.Background(),
		r.db,
		func(ctx context.Context) error {
			var err error

			for i, user := range entities {
				entities[i], err = r.Save(ctx, user)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)

	if saveError != nil {
		return []User{}, saveError
	}

	return entities, nil
}

func (r PostgresUserRepository) Count() int {
	var count int
	_ = r.db.Get(
		&count,
		fmt.Sprintf("SELECT COUNT(id) FROM %s", shema.Tables.User),
	)

	return count
}
func (r PostgresUserRepository) ExistsById(id int) bool {
	var exists bool
	_ = r.db.Get(
		&exists,
		fmt.Sprintf(
			"SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)",
			shema.Tables.User,
		),
		id,
	)

	return exists
}

func (r PostgresUserRepository) FindAll() ([]User, error) {
	users := make([]User, 0)

	err := r.db.Select(&users, fmt.Sprintf("SELECT * FROM %s", shema.Tables.User))
	if err != nil {
		return []User{}, nil
	}

	return users, nil
}

func (r PostgresUserRepository) FindById(id int) (User, error) {
	var user User
	err := r.db.Get(
		&user,
		fmt.Sprintf("SELECT * FROM %s WHERE id=$1", shema.Tables.User),
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf(
				"%w [id=%v]", ErrUserNotFound, id,
			)
		}

		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) FindAllById(ids []int) ([]User, error) { panic("not implemented") }

func (r PostgresUserRepository) FindByCredentials(username, password string) (User, error) {
	var user User

	err := r.db.Get(
		&user,
		fmt.Sprintf(
			"SELECT * FROM %s WHERE name=$1 AND password=$2",
			shema.Tables.User,
		),
		username, password,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, fmt.Errorf(
				"%w [username=%s, password=%s]",
				ErrUserNotFound, username, password,
			)
		}

		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) Delete(entity User) error         { panic("not implemented") }
func (r PostgresUserRepository) DeleteAll(entities ...User) error { panic("not implemented") }
func (r PostgresUserRepository) DeleteById(id int) error          { panic("not implemented") }
func (r PostgresUserRepository) DeleteAllById(ids []int) error    { panic("not implemented") }

func (r PostgresUserRepository) UpdatePassword(userId int, newPassword string) error {
	panic("not implemented")
}
func (r PostgresUserRepository) GrantRole(userId, roleId int) error { panic("not implemented") }

// func (r PostgresUserRepository) SaveUser(username, password, contact string) (int, error) {
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return -1, err
// 	}

// 	isUserExistQuery := fmt.Sprintf(`
// 		SELECT EXISTS (
// 			SELECT 1
// 			FROM public.%s
// 			WHERE name=$1
// 			LIMIT 1
// 		)`,
// 		shema.UserTable,
// 	)

// 	var isUserExist bool
// 	if err := tx.Get(&isUserExist, isUserExistQuery, username); err != nil {
// 		return -1, err
// 	}

// 	if isUserExist {
// 		return -1, fmt.Errorf(errUserAlreadyExists, username)
// 	}

// 	// Save user to the database
// 	createUserQuery := fmt.Sprintf(`
// 		INSERT INTO public.%s (
// 			name,
// 			password,
// 			contact
// 		)
// 		VALUES($1, $2, $3)
// 		RETURNING id`,
// 		shema.UserTable,
// 	)

// 	var userId int
// 	if err := tx.Get(&userId, createUserQuery, username, password, contact); err != nil {
// 		return -1, err
// 	}

// 	// Hardcode role to 'user'
// 	// if err := assignRole(tx, userId, 2); err != nil {
// 	// 	return -1, err
// 	// }

// 	// Commit transaction
// 	if err := tx.Commit(); err != nil {
// 		return -1, err
// 	}

// 	return userId, nil
// }

// func (r Postgres) GetUserByCredentials(username, password string) (*domain.User, error) {
// 	query := fmt.Sprintf(`
// 		SELECT *
// 		FROM public.%s
// 		WHERE name=$1 AND password=$2`,
// 		userTable,
// 	)

// 	var user domain.User
// 	if err := r.db.Get(&user, query, username, password); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, NewErrNotFound("user", "username, password",
// 				strings.Join([]string{username, password}, ", "))
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r Postgres) GetUserById(id int) (*domain.User, error) {
// 	getUserQuery := fmt.Sprintf(`
// 		SELECT
// 			u.id,
// 			u.name,
// 			u.password,
// 			u.contact,
// 			COALESCE (m.team_id, 0) AS team_id
// 		FROM public.%s AS u
// 		LEFT JOIN public.%s AS m ON u.id = m.user_id
// 		WHERE u.id=$1
// 		LIMIT 1`,
// 		userTable,
// 		teamMembersTable,
// 	)

// 	var user domain.User
// 	if err := r.db.Get(&user, getUserQuery, id); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, NewErrNotFound("user", "id", id)
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r Postgres) GetAllUsers() ([]*domain.User, error) {
// 	getUsersQuery := fmt.Sprintf(`
// 		SELECT
// 			u.id,
// 			u.name,
// 			u.password,
// 			u.contact,
// 			COALESCE (m.team_id, 0) AS team_id
// 		FROM public.%s AS u
// 		LEFT JOIN public.%s AS m ON u.id = m.user_id`,
// 		userTable,
// 		teamMembersTable)

// 	users := make([]*domain.User, 0)
// 	if err := r.db.Select(&users, getUsersQuery); err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// func (r Postgres) UpdatePassword(userId int, newPassword string) error {
// 	query := fmt.Sprintf(`
// 		UPDATE public.%s
// 		SET password = $1
// 		WHERE id = $2`,
// 		userTable,
// 	)

// 	if _, err := r.db.Exec(query, newPassword, userId); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r Postgres) AssignRole(userId, roleId int) error {
// 	// Begin transaction
// 	tx, err := r.db.Beginx()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	if err := assignRole(tx, userId, roleId); err != nil {
// 		return err
// 	}

// 	// Commit transaction
// 	if err := tx.Commit(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func assignRole(tx *sqlx.Tx, userId, roleId int) error {
// 	query := fmt.Sprintf(`
// 		INSERT INTO public.%s (
// 			user_id,
// 			role_id
// 		)
// 		VALUES($1, $2)`,
// 		userRolesTable)

// 	if _, err := tx.Exec(query, userId, roleId); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (r Postgres) GetUserRole(userId int) (*domain.Role, error) {
// 	query := fmt.Sprintf(`
// 		SELECT r.id, r.name
// 		FROM public.%s AS r
// 		INNER JOIN public.%s AS ur ON r.id=ur.role_id
// 		WHERE ur.user_id = $1;`,
// 		roleTable, userRolesTable)

// 	var role domain.Role
// 	if err := r.db.Get(&role, query, userId); err != nil {
// 		return nil, err
// 	}

// 	return &role, nil
// }

//func (r Postgres) IsUsernameExists(username string) error {
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
