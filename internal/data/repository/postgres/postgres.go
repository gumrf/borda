package postgres

import (
	"github.com/jmoiron/sqlx"

	"borda/internal/core/interfaces"
)

type PostgresRepository struct {
	db *sqlx.DB
}

// Verify interface compliance
var _ interfaces.Repository = (*PostgresRepository)(nil)

func (pr *PostgresRepository) Users() (interfaces.UserRepository, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Roles() (interfaces.RoleRepository, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Tasks() (interfaces.TaskRepository, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Teams() (interfaces.TeamRepository, error) {
	panic("not implemented") // TODO: implement me
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
