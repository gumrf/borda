package repository

import (
	"borda/internal/core"

	"github.com/jmoiron/sqlx"
)

// Repository contains all repositories
type Repository struct {
	db *sqlx.DB

	UserRepo     core.UserRepository
	TeamRepo     core.TeamRepository
	TaskRepo     core.TaskRepository
	SettingsRepo core.SettingsRepository
}

// Verify interface compliance
var _ core.Repository = (*Repository)(nil)

// NewRepository runs migrations. Returns new repository.
func NewRepository(db *sqlx.DB) (core.Repository, error) {
	repo := Repository{
		db:           db,
		UserRepo:     nil,
		TeamRepo:     nil,
		TaskRepo:     nil,
		SettingsRepo: nil,
	}

	return repo, nil
}

// Users returns UserRepository
func (r Repository) Users() core.UserRepository {
	return r.UserRepo
}

// Teams returns TeamRepository
func (r Repository) Teams() core.TeamRepository {
	return r.TeamRepo
}

// Tasks returns TaskRepository
func (r Repository) Tasks() core.TaskRepository {
	return r.TaskRepo
}

// Settings returns SettingsRepository
func (r Repository) Settings() core.SettingsRepository {
	return r.SettingsRepo
}
