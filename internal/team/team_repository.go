package team

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	// TODO: separate package
	shema "borda/internal/pkg"
	"borda/pkg/transaction"
)

type PostgresTeamRepository struct {
	db        *pgxpool.Pool
	txManager transaction.Manager
}

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrTeamNotFound      = errors.New("team not found")
	ErrTeamIsFull        = errors.New("team is full")
	ErrUserAlreadyInTeam = errors.New("user already in a team")
)

func NewTeamRepository(pool *pgxpool.Pool, txManager transaction.Manager) *PostgresTeamRepository {
	return &PostgresTeamRepository{
		db:        pool,
		txManager: txManager,
	}
}

func (r PostgresTeamRepository) Save(ctx context.Context, team Team) (Team, error) {
	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{
						"INSERT INTO %s (",
						"	name,",
						"	token",
						")",
						"VALUES($1, $2)",
						"RETURNING id",
					},
					"\n",
				),
				shema.Tables.Team,
			),
			team.Name, team.Token,
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
	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		for i, team := range entities {
			var err error
			entities[i], err = r.Save(ctx, team)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
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

func (r PostgresTeamRepository) FindAll(ctx context.Context) ([]Team, error) {
	var teams []Team

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		rows, err := tx.Query(
			ctx, fmt.Sprintf("SELECT * FROM %s", shema.Tables.Team),
		)
		if err != nil {
			return err
		}

		for rows.Next() {
			var team Team
			if err := rows.Scan(&team.Id, &team.Name, &team.Token); err != nil {
				return err
			}

			teams = append(teams, team)
		}

		return nil
	})

	// TODO:
	// Remove this code and make teams nil after scan error??
	if err != nil {
		return []Team{}, err
	}

	return teams, nil
}

func (r PostgresTeamRepository) find(ctx context.Context, fieldName string, fieldValue any) (Team, error) {
	var team Team

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{
						"SELECT *",
						"FROM %s",
						"WHERE %s=$1",
						"LIMIT 1",
					},
					"\n",
				),
				shema.Tables.Team, fieldName,
			),
			fieldValue,
		).Scan(&team.Id, &team.Name, &team.Token)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w [%s=%v]", ErrTeamNotFound, fieldName, fieldValue)
			}
		}

		return nil
	})

	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (r PostgresTeamRepository) FindAllById(ctx context.Context, ids []int) ([]Team, error) {
	teams := make([]Team, len(ids))
	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		for i, id := range ids {
			var err error
			team, err := r.FindById(ctx, id)
			if err != nil {
				return err
			}

			teams[i] = team
		}

		return nil
	})

	if err != nil {
		return []Team{}, err
	}

	return teams, nil
}

func (r PostgresTeamRepository) FindById(ctx context.Context, id int) (Team, error) {
	return r.find(ctx, "id", id)
}

func (r PostgresTeamRepository) FindByToken(ctx context.Context, token string) (Team, error) {
	return r.find(ctx, "token", token)
}

func (r PostgresTeamRepository) FindByUserId(ctx context.Context, userId int) (Team, error) {
	var team Team

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

		err := tx.QueryRow(
			ctx,
			fmt.Sprintf(
				strings.Join(
					[]string{
						"SELECT *",
						"FROM %s",
						"WHERE id=(",
						"	SELECT team_id",
						"	FROM team_member",
						"	WHERE user_id=$1",
						"	LIMIT 1",
						")",
						"LIMIT 1",
					},
					"\n",
				),
				shema.Tables.Team,
			),
			userId,
		).Scan(&team.Id, &team.Name, &team.Token)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("%w [userId=%v]", ErrTeamNotFound, userId)
			}
		}

		return nil
	})

	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (r PostgresTeamRepository) Delete(entity Team) error         { panic("not implemented") }
func (r PostgresTeamRepository) DeleteAll(entities ...Team) error { panic("not implemented") }
func (r PostgresTeamRepository) DeleteById(id int) error          { panic("not implemented") }
func (r PostgresTeamRepository) DeleteAllById(ids []int) error    { panic("not implemented") }

// TODO: Может лучше возвращать дополнительно ошибку??
func (r *PostgresTeamRepository) IsTeamFull(ctx context.Context, teamId int) bool {
	var isFull bool

	err := r.txManager.Run(ctx, func(ctx context.Context) error {
		tx, _ := r.txManager.GetTransaction(ctx)

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
	return r.txManager.Run(ctx, func(ctx context.Context) error {
		if r.IsTeamFull(ctx, teamId) {
			return fmt.Errorf("%w [id=%v]", ErrTeamIsFull, teamId)
		}

		tx, _ := r.txManager.GetTransaction(ctx)

		_, err := tx.Exec(
			ctx,
			fmt.Sprintf(
				"INSERT INTO %s (team_id, user_id, is_captain) VALUES($1, $2, $3)",
				shema.Tables.TeamMember,
			),
			teamId, userId, false,
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
}
