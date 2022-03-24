package postgres

import (
	"borda/internal/domain"
	"borda/internal/repository"
	"errors"

	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r TeamRepository) SaveTeam(teamLeaderId int, teamName string) (int, error) {
	// Begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, err
	}

	isTeamExistQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		teamTable,
	)

	// Check if team name already exists in database
	var isTeamExists bool
	if err := tx.Get(&isTeamExists, isTeamExistQuery, teamName); err != nil {
		return -1, err
	}

	if isTeamExists {
		return -1, repository.NewErrNotFound("team", "name", teamName)
	}

	// Generate access token for team
	uuid := uuid.New().String()

	// Save team to database
	saveTeamQuery := fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			token,
			team_leader_id
		) 
		VALUES($1, $2, $3)
		RETURNING id`,
		teamTable,
	)

	var id int
	row := tx.QueryRow(saveTeamQuery, teamName, uuid, teamLeaderId)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil
}

func (r TeamRepository) GetTeamById(teamId int) (*domain.Team, error) {
	getTeamQuery := fmt.Sprintf(`
		SELECT * 
		FROM public.%s 
		WHERE id=$1
		LIMIT 1`,
		teamTable,
	)

	var team domain.Team
	if err := r.db.Get(&team, getTeamQuery, teamId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewErrNotFound("team", "id", teamId)
		}
		return nil, err
	}

	return &team, nil
}

func (r TeamRepository) GetTeamByToken(token string) (*domain.Team, error) {
	getTeamQuery := fmt.Sprintf(`
		SELECT * 
		FROM public.%s 
		WHERE token=$1
		LIMIT 1`,
		teamTable,
	)

	var team domain.Team
	if err := r.db.Get(&team, getTeamQuery, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewErrNotFound("team", "token", token)
		}
		return nil, err
	}

	return &team, nil
}

func (r TeamRepository) AddMember(teamId, userId int) error {
	// Begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	// Check if user id exists in database
	isUserExistQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE id=$1
			LIMIT 1
		)`,
		userTable,
	)

	var isUserExist bool
	if err := tx.Get(&isUserExist, isUserExistQuery, userId); err != nil {
		return err
	}

	if !isUserExist {
		return repository.NewErrNotFound("user", "id", userId)
	}

	// Check if team id exists in database
	isTeamExistQuery := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE id=$1
			LIMIT 1
		)`,
		teamTable,
	)

	var isTeamExist bool
	if err := tx.Get(&isTeamExist, isTeamExistQuery, teamId); err != nil {
		return err
	}

	if !isTeamExist {
		return repository.NewErrNotFound("team", "id", teamId)
	}

	// Get the number of members in the team
	var teamMembersCount int
	teamMembersCountQuery := fmt.Sprintf(`
		SELECT COUNT(user_id)
		FROM %s
		WHERE team_id=$1`,
		teamMembersTable,
	)

	if err := tx.Get(&teamMembersCount, teamMembersCountQuery, teamId); err != nil {
		return err
	}

	// Get team members limit from settings
	var teamMembersLimit string
	teamMembersLimitQuery := fmt.Sprintf(`
		SELECT value 
		FROM %s
		WHERE key=$1`,
		settingsTable,
	)
	if err := tx.Get(&teamMembersLimit, teamMembersLimitQuery, "team_limit"); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.NewErrNotFound("setting", "value", "team_limit")
		}
		return err
	}

	// Convert team limit setting value from string to int
	teamMembersLimitInt, err := strconv.Atoi(teamMembersLimit)
	if err != nil {
		return err
	}

	if teamMembersCount >= teamMembersLimitInt {
		return errors.New("team is full")
	}

	// Attach user to the team
	addMemberQuery := fmt.Sprintf(`
		INSERT INTO public.%s (
			team_id, 
			user_id
		) 
		VALUES($1, $2)
		RETURNING id`,
		teamMembersTable,
	)

	var id int = -1
	if err = tx.Get(&id, addMemberQuery, teamId, userId); err != nil || id == -1 {
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//func (r TeamRepository) GetMembers(teamId int) (users []domain.User, err error) {
//	// Check team exist
//	query := fmt.Sprintf(`
//		SELECT id
//		FROM %s
//		WHERE id=$1`,
//		r.teamMembersTable,
//	)
//	var team_id int
//	err = r.db.QueryRowx(query, teamId).Scan(&team_id)
//	if err != nil {
//		return []domain.User{}, fmt.Errorf("team repository getMembers error: Team not found with id=%v", teamId)
//	}
//
//	// Get
//	query = fmt.Sprintf(`
//		SELECT *
//		FROM %s
//		WHERE ID IN (
//			SELECT user_id
//			FROM %s
//			WHERE team_id=$1
//		)`,
//		r.userTable,
//		r.teamMembersTable,
//	)
//
//	var _users = make([]domain.User, 0)
//
//	rows, err := r.db.Queryx(query, teamId)
//	if err != nil {
//		return []domain.User{}, fmt.Errorf("team repository getMembers error: Members not found in team with id=%v, %v", teamId, err)
//	}
//
//	for rows.Next() {
//		var user domain.User
//		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Contact)
//		if err != nil {
//			return []domain.User{}, fmt.Errorf("team repository getMembers error: On convert to domain in team with id=%v, %v", teamId, err)
//		}
//
//		// user.TeamId = teamId
//		_users = append(_users, user)
//	}
//
//	return _users, nil
//}
