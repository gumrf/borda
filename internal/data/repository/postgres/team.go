package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresTeamRepository struct {
	db *sqlx.DB
	tableName string
}

var _ interfaces.TeamRepository = (*PostgresTeamRepository)(nil)

func NewPostgresTeamRepository(db *sqlx.DB) interfaces.TeamRepository {
	return PostgresTeamRepository{db: db, tableName: "team"}
}

func (r PostgresTeamRepository) Create(teamLeaderId int, teamName string) (team entity.Team, err error) {
	fmt.Printf("\n\n%s\n\n", r.tableName)
	// query := "INSERT INTO user (name, password, contact) VALUES(?, ?, ?)"

	// result, err := r.db.Exec(qwery, username, password, contact)
	// if err != nil {
	// 	return -1, err
	// }

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return -1, err
	// }

	// return int(id), nil
	return entity.Team{}, nil
}
func (r PostgresTeamRepository) AddMember(teamId, userId int) error {
	return nil
}
func (r PostgresTeamRepository) Get(teamId int) (team entity.Team, err error) {
	return entity.Team{}, nil
}
