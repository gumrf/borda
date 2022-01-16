package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostgresTeamRepository struct {
	db                   *sqlx.DB
	tableTeamName        string
	tableUserName        string
	tableTeamMembersName string
}

var _ interfaces.TeamRepository = (*PostgresTeamRepository)(nil)

func NewPostgresTeamRepository(db *sqlx.DB) interfaces.TeamRepository {
	return PostgresTeamRepository{db: db, tableTeamName: "team", tableUserName: "\"user\"", tableTeamMembersName: "team_members"}
}

func (r PostgresTeamRepository) Create(teamLeaderId int, teamName string) (team entity.Team, err error) {
	// Generate uuid
	uuid := uuid.New().String()

	// Write db
	query := fmt.Sprintf(
		`INSERT INTO public.%s 
		(name, token, team_leader_id) 
		VALUES($1, $2, $3)
		RETURNING id`,
		r.tableTeamName,
	)

	id := -1
	err = r.db.QueryRow(query, teamName, uuid, teamLeaderId).Scan(&id)
	if err != nil || id == -1 {
		_err := fmt.Errorf("team repository create error: %v", err)
		return entity.Team{}, _err
	}

	// Build obj
	obj := entity.Team{
		Id:           id,
		Name:         teamName,
		TeamLeaderId: teamLeaderId,
		Token:        uuid,
	}
	return obj, nil
}

func (r PostgresTeamRepository) AddMember(teamId, userId int) error {
	// Query check result struct
	type QResult struct {
		TeamId        sql.NullInt64 `db:"team_id"`
		UserId        sql.NullInt64 `db:"user_id"`
		TeamMembersId sql.NullInt64 `db:"tm_id"`
	}

	// Get team_id, user_id, team_members_id for check
	// Select like this team_id | user_id | team_members_id
	query := fmt.Sprintf(
		`SELECT COALESCE((
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
		_err := fmt.Errorf("team repository addMember error: %v", err)
		return _err
	}

	t_id, err := result.TeamId.Value()
	if t_id == nil || err != nil {
		_err := fmt.Errorf("team repository addMember error: Team with id=%v not found", teamId)
		return _err
	}

	u_id, err := result.UserId.Value()
	if u_id == nil || err != nil {
		_err := fmt.Errorf("team repository addMember error: User with id=%v not found", userId)
		return _err
	}

	tm, err := result.TeamMembersId.Value()
	if tm != nil {
		_err := fmt.Errorf("team repository addMember error: User id=%v already in team with id=%v", userId, teamId)
		return _err
	}
	if err != nil {
		_err := fmt.Errorf("team repository addMember error: %v", err)
		return _err
	}

	// Write db
	query = fmt.Sprintf(
		`INSERT INTO public.%s 
		(team_id, user_id) 
		VALUES($1, $2)
		RETURNING id`,
		r.tableTeamMembersName,
	)

	id := -1
	err = r.db.QueryRow(query, teamId, userId).Scan(&id)
	if err != nil || id == -1 {
		_err := fmt.Errorf("team repository addMember error: %v", err)
		return _err
	}
	return nil
}

func (r PostgresTeamRepository) Get(teamId int) (team entity.Team, err error) {
	// Get
	query := fmt.Sprintf(
		`SELECT * FROM public.%s 
		WHERE id=$1`,
		r.tableTeamName,
	)

	// Build empty obj
	obj := entity.Team{
		Id:           -1,
		Name:         "",
		TeamLeaderId: -1,
		Token:        "",
	}

	// Scan to struct, fill obj
	err = r.db.QueryRowx(query, teamId).StructScan(&obj)
	
	if err == sql.ErrNoRows {
		_err := fmt.Errorf("team repository get error: Team not found with id=%v", teamId)
		return entity.Team{}, _err
	}
	if err != nil {
		_err := fmt.Errorf("team repository get error: %v", err)
		return entity.Team{}, _err
	}

	// Check
	if obj.Id == -1 || obj.TeamLeaderId == -1 || obj.Name == "" || obj.Token == "" {
		_err := fmt.Errorf("team repository get error: Field of struct not filled")
		return entity.Team{}, _err
	}

	return obj, nil
}
