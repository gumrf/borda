package team

import (
	// TODO: pkg.Tables
	shema "borda/internal/pkg"
	"borda/pkg/transaction"
	"context"
	"errors"
	"strings"

	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

type PostgresTeamRepository struct {
	db *pgxpool.Pool
}

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamIsFull        = errors.New("team is full")
	ErrUserAlreadyInTeam = errors.New("user already in a team")
)

func NewTeamRepository(pool *pgxpool.Pool) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: pool}
}

func (r PostgresTeamRepository) Save(ctx context.Context, team Team) (Team, error) {
	err := transaction.Transactional(ctx, r.db, func(ctx context.Context) error {
		var tx pgx.Tx

		tx, ok := transaction.ExtractTransaction(ctx)
		if !ok {
			panic("this is not possible")
		}

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{
						"INSERT INTO %s (",
						"	name,",
						"	token,",
						"	team_leader_id",
						")",
						"VALUES($1, $2, $3)",
						"RETURNING id",
					},
					"\n",
				),
				shema.Tables.Team,
			),
			team.Name, team.Token, team.TeamLeaderId,
		).Scan(&team.Id)
		if err != nil {
			var pgxError *pgconn.PgError
			if errors.As(err, &pgxError) {
				switch pgxError.Code {
				case "23503": //unique_violation
					return fmt.Errorf(
						"%w [name=%v]", ErrTeamAlreadyExists, team.Name,
					)
				}
			}

			return fmt.Errorf("scan: %w", err)
		}

		return nil
	})

	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (r *PostgresTeamRepository) SaveAll(ctx context.Context, entities []Team) ([]Team, error) {
	err := transaction.Transactional(ctx, r.db, func(ctx context.Context) error {
		for i, team := range entities{
			var err error
			entities[i], err = r.Save(ctx, team)
			if err != nil{
				return err
			}
		}

		return nil
	})
	
	if err != nil{
		return []Team{}, err
	}

	return entities, nil
}

func (r PostgresTeamRepository) Count() int {
	var count int
	_ = r.db.QueryRow(
		context.TODO(),
		fmt.Sprintf("SELECT COUNT(id) FROM %s", shema.Tables.Team),
	).Scan(&count)

	return count
}

func (r PostgresTeamRepository) ExistsById(id int) bool {
	var exists bool
	_ = r.db.QueryRow(
		context.TODO(),
		fmt.Sprintf(
			"SELECT EXISTS(SELECT 1 FROM %s WHERE id=$1)", shema.Tables.Team,
		),
		id,
	).Scan(&exists)

	return exists
}

func (r PostgresTeamRepository) FindAll() ([]Team, error) {
	// tx, err := r.db.Beginx()
	// if err != nil {
	// 	return nil, err
	// }
	// defer tx.Rollback() // nolint

	// getTeamsQuery := fmt.Sprintf(`
	// 	SELECT *
	// 	FROM public.%s`,
	// 	shema.TeamTable,
	// )

	// teams := make([]*domain.Team, 0)
	// if err := tx.Select(&teams, getTeamsQuery); err != nil {
	// 	return nil, err
	// }

	// for _, team := range teams {
	// 	getMembersQuery := fmt.Sprintf(`
	// 	SELECT name AS user_name, id AS user_id
	// 	FROM %s
	// 	WHERE ID IN (
	// 		SELECT user_id
	// 		FROM %s
	// 		WHERE team_id=$1
	// 	)`,
	// 		shema.UserTable,
	// 		shema.TeamMembersTable,
	// 	)

	// 	if err := tx.Select(&team.Members, getMembersQuery, team.Id); err != nil {
	// 		return nil, err
	// 	}

	// }

	// return teams, nil
	panic("not implemented")
}
func (r PostgresTeamRepository) FindAllById(ids []int) ([]Team, error) {
	panic("not implemented")
}

func (r PostgresTeamRepository) FindById(id int) (Team, error) {
	// var team Team
	// getTeamQuery := fmt.Sprintf(
	// 	`SELECT * FROM %s WHERE id=$1 LIMIT 1`,
	// 	shema.TeamTable,
	// )

	// if err := r.db.Get(&team, getTeamQuery, teamId); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return Team{}, fmt.Errorf("team with id %s not found", teamId)
	// 	}
	// 	return Team{}, err
	// }

	// 	err := r.db.Select(
	// 		&team.Members,
	// 		fmt.Sprintf(`
	// SELECT name AS user_name, id AS user_id
	// FROM %s
	// WHERE ID IN (
	// 	SELECT user_id
	// 	FROM %s
	// 	WHERE team_id=$1
	// )`,
	// 			shema.Tables.User,
	// 			shema.Tables.TeamMember,
	// 		),
	// 		id,
	// 	)
	// 	if err != nil {
	// 		return Team{}, err
	// 	}

	// 	return team, nil
	panic("not implemented")
}

func (r PostgresTeamRepository) FindByToken(token string) (Team, error) {
	// getTeamQuery := fmt.Sprintf(`
	// 	SELECT *
	// 	FROM public.%s
	// 	WHERE token=$1
	// 	LIMIT 1`,
	// 	shema.Tables.Team,
	// )

	// var team Team
	// if err := r.db.Get(&team, getTeamQuery, token); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return Team{}, fmt.Errorf("%w [token=%s", ErrTeamNotFound, token)
	// 	}
	// 	return Team{}, err
	// }

	// return team, nil
	panic("not implemented")
}

func (r PostgresTeamRepository) FindByUserId(userId int) (Team, error) {
	// getTeamQuery := fmt.Sprintf(`
	// 	SELECT team_id
	// 	FROM public.%s
	// 	WHERE user_id=$1
	// 	LIMIT 1`,
	// 	shema.Tables.TeamMember,
	// )

	// var team Team
	// if err := r.db.Get(&team, getTeamQuery, userId); err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		return Team{}, fmt.Errorf("%w [userId=%v]", ErrTeamNotFound, userId)
	// 	}
	// 	return Team{}, err
	// }

	// return team, nil
	panic("not implemented")
}

func (r PostgresTeamRepository) Delete(entity Team) error         { panic("not implemented") }
func (r PostgresTeamRepository) DeleteAll(entities ...Team) error { panic("not implemented") }
func (r PostgresTeamRepository) DeleteById(id int) error          { panic("not implemented") }
func (r PostgresTeamRepository) DeleteAllById(ids []int) error    { panic("not implemented") }

// TODO: Может лучше возвращать дополнительно ошибку??
func (r *PostgresTeamRepository) IsTeamFull(ctx context.Context, teamId int) bool {
	var isFull bool

	err := transaction.Transactional(ctx, r.db, func(ctx context.Context) error {
		tx, _ := transaction.ExtractTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{
						"SELECT(",
						"	(SELECT COUNT(user_id) FROM %s WHERE team_id=$1)",
						"	>=",
						"	(SELECT CAST(value AS INTEGER) FROM %s WHERE key='team_limit')",
						") as is_full",
					},
					"\n",
				),
				shema.Tables.TeamMember, shema.Tables.Settings,
			),
			teamId,
		).Scan(&isFull)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return false
	}

	return isFull
}

func (r PostgresTeamRepository) AddMember(ctx context.Context, teamId, userId int) error {
	err := transaction.Transactional(ctx, r.db, func(ctx context.Context) error {
		if r.IsTeamFull(ctx, teamId) {
			return fmt.Errorf("%w [id=%v]", ErrTeamIsFull, teamId)
		}

		tx, _ := transaction.ExtractTransaction(ctx)

		_, err := tx.Exec(
			ctx,
			fmt.Sprintf(
				"INSERT INTO %s (team_id, user_id) VALUES($1, $2)",
				shema.Tables.TeamMember,
			),
			teamId, userId,
		)
		if err != nil {
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr) {
				switch pgerr.ConstraintName {
				// TODO: May be not all errors caught???
				case "team_member_user_id_key": // 23505 - duplicate key value violates unique constraint
					return fmt.Errorf("%w [userId=%v]", ErrUserAlreadyInTeam, userId)
				}
			}

			return err
		}

		return nil
	})

	return err
}

func (r PostgresTeamRepository) GetMembers(teamId int) (users []User, err error) {
	// Check team exist
	// query := fmt.Sprintf(`
	// 	SELECT id
	// 	FROM %s
	// 	WHERE id=$1`,
	// 	shema.Tables.TeamMember,
	// )
	// var team_id int
	// err = r.db.QueryRowx(query, teamId).Scan(&team_id)
	// if err != nil {
	// 	return []User{}, fmt.Errorf("team with id %v not found", team_id)
	// }

	// // Get
	// query = fmt.Sprintf(`
	// 	SELECT *
	// 	FROM %s
	// 	WHERE ID IN (
	// 		SELECT user_id
	// 		FROM %s
	// 		WHERE team_id=$1
	// 	)`,
	// 	shema.Tables.User,
	// 	shema.Tables.TeamMember,
	// )

	// var _users = make([]User, 0)

	// rows, err := r.db.Queryx(query, teamId)
	// if err != nil {
	// 	return []User{}, fmt.Errorf("member with id %v not found", teamId)
	// }

	// for rows.Next() {
	// 	var user User
	// 	err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Contact)
	// 	if err != nil {
	// 		return []User{}, fmt.Errorf("team repository getMembers error: On convert to domain in team with id=%v, %v", teamId, err)
	// 	}

	// 	// user.TeamId = teamId
	// 	_users = append(_users, user)
	// }

	// return _users, nil
	panic("not implemented")
}
