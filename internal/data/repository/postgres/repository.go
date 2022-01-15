package repository

import (
	"github.com/jmoiron/sqlx"

	"borda/internal/core/interfaces"
)

type PostgresRepository struct {
	db    *sqlx.DB
	repos repositories
}

type repositories struct {
	user interfaces.UserRepository
	team interfaces.TeamRepository
	task interfaces.TaskRepository
}

// Verify interface compliance
var _ interfaces.Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *sqlx.DB) PostgresRepository {
	repository := PostgresRepository{db: db}

	repository.repos.user = newPostgresUserRepository(db)
	repository.repos.team = newPostgresTeamRepository(db)
	repository.repos.task = newPostgresTaskRepository(db)

	return repository
}

func (r *PostgresRepository) Users() interfaces.UserRepository {
	return r.repos.user
}

func (r *PostgresRepository) Teams() interfaces.TeamRepository {
	return r.repos.team
}

func (r *PostgresRepository) Tasks() interfaces.TaskRepository {
	return r.repos.task
}
