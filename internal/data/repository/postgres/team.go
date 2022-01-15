package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresTeamRepository struct {
	db *sqlx.DB
}

var _ interfaces.TeamRepository = (*PostgresTeamRepository)(nil)

func newPostgresTeamRepository(db *sqlx.DB) PostgresTeamRepository {
	return PostgresTeamRepository{db: db}
}

func (r PostgresTeamRepository) Create(teamLeaderId int, teamName string) (team entity.Team, err error)
func (r PostgresTeamRepository) AddMember(teamId, userId int) error
func (r PostgresTeamRepository) Get(teamId int) (team entity.Team, err error)
