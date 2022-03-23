package postgres

import (
	"borda/internal/domain"

	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db                   *sqlx.DB
	tableTeamName        string
	tableUserName        string
	tableTeamMembersName string
	tableSettingsName    string
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{
		db:                   db,
		tableTeamName:        "team",
		tableUserName:        "\"user\"",
		tableTeamMembersName: "team_member",
		tableSettingsName:    "settings",
	}
}

func (r TeamRepository) CreateNewTeam(teamLeaderId int, teamName string) (int, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		r.tableTeamName)

	var isTeamNameExists bool
	err := r.db.Get(&isTeamNameExists, query, teamName)
	if err != nil {
		return -1, err
	}

	if isTeamNameExists {
		return -1, fmt.Errorf("team name already exists")
	}

	// Generate uuid
	uuid := uuid.New().String()

	query = fmt.Sprintf(`
		INSERT INTO public.%s (
			name,
			token,
			team_leader_id
		) 
		VALUES($1, $2, $3)
		RETURNING id`,
		r.tableTeamName,
	)

	var id int
	row := r.db.QueryRow(query, teamName, uuid, teamLeaderId)
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("TeamRepository.Create: %w", err)
	}

	return id, nil
}

func (r TeamRepository) IsTeamNameExists(teamName string) error {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE name=$1
			LIMIT 1
		)`,
		r.tableTeamName)

	var isTeamNameExists bool
	err := r.db.Get(&isTeamNameExists, query, teamName)
	if err != nil {
		return err
	}

	if isTeamNameExists {
		return domain.ErrTeamAlreadyExists
	}

	return nil
}

func (r TeamRepository) GetTeamById(teamId int) (domain.Team, error) {
	query := fmt.Sprintf(`
		SELECT * 
		FROM public.%s 
		WHERE id=$1
		LIMIT 1`,
		r.tableTeamName,
	)

	var team domain.Team
	err := r.db.Get(&team, query, teamId)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (r TeamRepository) GetTeamByToken(token string) (domain.Team, error) {
	query := fmt.Sprintf(`
		SELECT * 
		FROM public.%s 
		WHERE token=$1
		LIMIT 1`,
		r.tableTeamName,
	)

	var team domain.Team
	err := r.db.Get(&team, query, token)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (r TeamRepository) IsTeamTokenExists(token string) error {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM public.%s
			WHERE token=$1
			LIMIT 1
		)`,
		r.tableTeamName)

	var isTeamTokenValid bool
	err := r.db.Get(&isTeamTokenValid, query, token)
	if err != nil {
		return err
	}

	if !isTeamTokenValid {
		return domain.ErrTeamTokenIsInvalid
	}

	return nil
}

func (r TeamRepository) AddMember(teamId, userId int) error {
	// Query check result struct
	type QResult struct {
		TeamId        sql.NullInt64 `db:"team_id"`
		UserId        sql.NullInt64 `db:"user_id"`
		TeamMembersId sql.NullInt64 `db:"tm_id"`
	}

	// Get team_id, user_id, team_members_id for check
	// Select like this team_id | user_id | team_members_id
	query := fmt.Sprintf(`
		SELECT COALESCE((
			SELECT id FROM public.%s
			WHERE id=$1), NULL
		) as team_id, 
		COALESCE((
			SELECT id FROM public.%s
			WHERE id=$2), NULL
		) as user_id, 
		COALESCE((
			SELECT id FROM public.%s
			WHERE team_id=$1 AND user_id=$2), NULL
		) as tm_id`,
		r.tableTeamName,
		r.tableUserName,
		r.tableTeamMembersName,
	)

	result := QResult{
		TeamId:        sql.NullInt64{},
		UserId:        sql.NullInt64{},
		TeamMembersId: sql.NullInt64{},
	}

	// Scan to struct, fill obj
	err := r.db.QueryRowx(query, teamId, userId).StructScan(&result)
	if err != nil {
		return fmt.Errorf("team repository addMember error: %v", err)
	}

	t_id, err := result.TeamId.Value()
	if t_id == nil || err != nil {
		return fmt.Errorf("team repository addMember error: Team with id=%v not found", teamId)
	}

	u_id, err := result.UserId.Value()
	if u_id == nil || err != nil {
		return fmt.Errorf("team repository addMember error: User with id=%v not found", userId)
	}

	tm, err := result.TeamMembersId.Value()
	if tm != nil {
		return fmt.Errorf("team repository addMember error: User id=%v already in team with id=%v", userId, teamId)
	}
	if err != nil {
		return fmt.Errorf("team repository addMember error: %v", err)
	}

	// Check limit
	// Tested manual, it really works, trust me :)
	var valueLimit string
	query = fmt.Sprintf(`
		SELECT value 
		FROM %s
		WHERE key=$1`,
		r.tableSettingsName,
	)
	err = r.db.QueryRowx(query, "team_limit").Scan(&valueLimit)
	if err != nil {
		return fmt.Errorf("team repository addMember error: Not found team_limit in db, %v", err)
	}

	memberLimit, err := strconv.Atoi(valueLimit)
	if err != nil {
		return fmt.Errorf("team repository addMember error: team_limit in db not converted to integer, %v", err)
	}
	var alreadyExistMembers int
	query = fmt.Sprintf(`
		SELECT COUNT(user_id)
		FROM %s
		WHERE team_id=$1`,
		r.tableTeamMembersName,
	)
	err = r.db.QueryRow(query, teamId).Scan(&alreadyExistMembers)
	if err != nil {
		return fmt.Errorf("team repository addMember error: %v", err)
	}
	if alreadyExistMembers+1 > memberLimit {
		return fmt.Errorf("team repository addMember Limit team error. Already members: %v, limit: %v", alreadyExistMembers, memberLimit)
	}

	// Write db
	query = fmt.Sprintf(`
		INSERT INTO public.%s (
			team_id, 
			user_id
		) 
		VALUES($1, $2)
		RETURNING id`,
		r.tableTeamMembersName,
	)

	id := -1
	err = r.db.QueryRow(query, teamId, userId).Scan(&id)
	if err != nil || id == -1 {
		return fmt.Errorf("team repository addMember error: %v", err)
	}
	return nil
}

func (r TeamRepository) IsTeamFull(teamId int) error {
	var valueLimit string
	query := fmt.Sprintf(`
		SELECT value 
		FROM %s
		WHERE key=$1`,
		r.tableSettingsName,
	)
	err := r.db.QueryRowx(query, "team_limit").Scan(&valueLimit)
	if err != nil {
		return fmt.Errorf("team repository addMember error: Not found team_limit in db, %v", err)
	}

	memberLimit, err := strconv.Atoi(valueLimit)
	if err != nil {
		return fmt.Errorf("team repository addMember error: team_limit in db not converted to integer, %v", err)
	}

	var alreadyExistMembers int
	query = fmt.Sprintf(`
		SELECT COUNT(user_id)
		FROM %s
		WHERE team_id=$1`,
		r.tableTeamMembersName,
	)
	err = r.db.QueryRow(query, teamId).Scan(&alreadyExistMembers)
	if err != nil {
		return fmt.Errorf("team repository addMember error: %v", err)
	}
	if alreadyExistMembers+1 > memberLimit {
		return fmt.Errorf("team repository addMember Limit team error. Already members: %v, limit: %v", alreadyExistMembers, memberLimit)
	}

	return nil
}

func (r TeamRepository) GetMembers(teamId int) (users []domain.User, err error) {
	// Check team exist
	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE id=$1`,
		r.tableTeamMembersName,
	)
	var team_id int
	err = r.db.QueryRowx(query, teamId).Scan(&team_id)
	if err != nil {
		return []domain.User{}, fmt.Errorf("team repository getMembers error: Team not found with id=%v", teamId)
	}

	// Get
	query = fmt.Sprintf(`
		SELECT *
		FROM %s 
		WHERE ID IN (
			SELECT user_id 
			FROM %s
			WHERE team_id=$1
		)`,
		r.tableUserName,
		r.tableTeamMembersName,
	)

	var _users = make([]domain.User, 0)

	rows, err := r.db.Queryx(query, teamId)
	if err != nil {
		return []domain.User{}, fmt.Errorf("team repository getMembers error: Members not found in team with id=%v, %v", teamId, err)
	}

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Contact)
		if err != nil {
			return []domain.User{}, fmt.Errorf("team repository getMembers error: On convert to domain in team with id=%v, %v", teamId, err)
		}

		// user.TeamId = teamId
		_users = append(_users, user)
	}

	return _users, nil
}
