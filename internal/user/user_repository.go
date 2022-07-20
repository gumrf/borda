package user

import (
	shema "borda/internal/pkg"
	"borda/pkg/transaction"
	"strings"

	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
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
	db        *pgxpool.Pool
	txManager transaction.Manager
}

func NewUserRepository(pool *pgxpool.Pool, txManager transaction.Manager) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:        pool,
		txManager: txManager,
	}
}

func (r PostgresUserRepository) Save(ctx context.Context, user User) (User, error) {
	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		// query := fmt.Sprintf(`
		// INSERT INTO public.%s (
		// 	name,
		// 	password,
		// 	contact
		// )
		// VALUES($1, $2, $3)
		// RETURNING id`, shema.Tables.User,
		// )

		// `
		// 		INSERT INTO %s (
		// 			name,
		// 			password,
		// 			contact
		// 		)
		// 		VALUES($1, $2, $3)
		// 		RETURNING id`, shema.Tables.User,

		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{"INSERT INTO %s (",
						" name,",
						" password,",
						" contact",
						")",
						"VALUES ($1, $2, $3)",
						"RETURNING id",
					},
					"\n",
				),
				shema.Tables.User,
			),
			user.Username, user.Password, user.Contact,
		).Scan(&user.Id)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) SaveAll(ctx context.Context, entities []User) ([]User, error) {
	saveError := r.txManager.Run(ctx, func(ctx context.Context) error {
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
	_ = r.db.QueryRow(
		context.TODO(),
		fmt.Sprintf("SELECT COUNT(id) FROM %s", shema.Tables.User),
	).Scan(&count)

	return count
}
func (r PostgresUserRepository) ExistsById(id int) bool {
	var exists bool
	_ = r.db.QueryRow(
		context.TODO(),
		fmt.Sprintf(
			"SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)",
			shema.Tables.User,
		),
		id,
	).Scan(&exists)

	return exists
}

func (r PostgresUserRepository) FindAll(ctx context.Context) ([]User, error) {
	users := make([]User, 0)

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		rows, err := tx.Query(
			ctx, fmt.Sprintf("SELECT * FROM %s", shema.Tables.User),
		)
		if err != nil {
			return nil
		}

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.Id, &user.Username,
				&user.Password, &user.Contact, &user.TeamId); err != nil {
				return err
			}

			users = append(users, user)
		}

		return nil
	})

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (r PostgresUserRepository) FindById(ctx context.Context, id int) (User, error) {
	var user User

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf("SELECT * FROM %s WHERE id=$1", shema.Tables.User),
			id,
		).Scan(&user.Id, &user.Username, &user.Password,
			&user.Contact, &user.TeamId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf(
					"%w [id=%v]", ErrUserNotFound, id,
				)
			}
			return err
		}

		return nil
	})

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) FindAllById(ctx context.Context, ids []int) ([]User, error) {
	users := make([]User, len(ids))

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		for _, id := range ids {
			var err error
			user, err := r.FindById(ctx, id)
			if err != nil {
				return err
			}
			users = append(users, user)
		}
		return nil
	})

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

func (r PostgresUserRepository) FindByCredentials(ctx context.Context, username, password string) (User, error) {
	var user User

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				"SELECT * FROM %s WHERE name=$1 AND password=$2",
				shema.Tables.User,
			),
			username, password,
		).Scan(&user)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf(
					"%w [username=%s, password=%s]",
					ErrUserNotFound, username, password,
				)
			}

			return err
		}

		return nil
	})

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r PostgresUserRepository) Delete(entity User) error         { panic("not implemented") }
func (r PostgresUserRepository) DeleteAll(entities ...User) error { panic("not implemented") }
func (r PostgresUserRepository) DeleteById(id int) error          { panic("not implemented") }
func (r PostgresUserRepository) DeleteAllById(ids []int) error    { panic("not implemented") }

func (r PostgresUserRepository) UpdatePassword(ctx context.Context, userId int, newPassword string) (User, error) {
	// 	query := fmt.Sprintf(`
	// 		UPDATE public.%s
	// 		SET password = $1
	// 		WHERE id = $2`,
	// 		userTable,
	// 	)
	var user User

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(`
			 		UPDATE %s SET password = $1 WHERE id = $2`,
				shema.Tables.User),
			newPassword, userId,
		).Scan(&user.Id, &user.Username, &user.Password, &user.Contact, &user.TeamId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf(
					"%w [user_id=%v, password=%s]",
					ErrUserNotFound, userId, newPassword,
				)
			}

			return err
		}

		return nil
	})

	if err != nil {
		return User{}, err
	}

	return user, nil
}
func (r PostgresUserRepository) GrantRole(ctx context.Context, userId, roleId int) error {
	panic("not implemented")
}

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
