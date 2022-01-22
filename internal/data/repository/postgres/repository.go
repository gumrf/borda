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
	user     interfaces.UserRepository
	team     interfaces.TeamRepository
	task     interfaces.TaskRepository
	settings interfaces.SettingsRepository
}

// Verify interface compliance
var _ interfaces.Repository = (*PostgresRepository)(nil)

func NewPostgresRepository(db *sqlx.DB) interfaces.Repository {
	repository := PostgresRepository{db: db}

	repository.repos.user = NewPostgresUserRepository(db)
	repository.repos.team = NewPostgresTeamRepository(db)
	repository.repos.task = NewPostgresTaskRepository(db)
	repository.repos.settings = NewPostgresSettingsRepository(db)

	return repository
}

func (r PostgresRepository) Users() interfaces.UserRepository {
	return r.repos.user
}

func (r PostgresRepository) Settings() interfaces.SettingsRepository {
	return r.repos.settings
}

func (r PostgresRepository) Teams() interfaces.TeamRepository {
	return r.repos.team
}

func (r PostgresRepository) Tasks() interfaces.TaskRepository {
	return r.repos.task
}
