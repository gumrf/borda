package postgres

import (
	"github.com/jmoiron/sqlx"

	"borda/internal/repository"
)

type PostgresRepository struct {
	db *sqlx.DB
}

// Verify interface compliance
var _ repository.RepositoryI = (*PostgresRepository)(nil)

func (pr *PostgresRepository) Users() (repository.UserRepositoryI, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Roles() (repository.RoleRepositoryI, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Tasks() (repository.TaskRepositoryI, error) {
	panic("not implemented") // TODO: implement me
}

func (pr *PostgresRepository) Teams() (repository.TeamRepositoryI, error) {
	panic("not implemented") // TODO: implement me
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}
