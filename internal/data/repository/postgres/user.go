package repository

import (
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresUserRepository struct {
	db *sqlx.DB
}

var _ interfaces.UserRepository = (*PostgresUserRepository)(nil)

func newPostgresUserRepository(db *sqlx.DB) PostgresUserRepository {
	return PostgresUserRepository{db: db}
}

func (r PostgresUserRepository) Create(username, password, contact string) (userId int, err error)
func (r PostgresUserRepository) UpdatePassword(userId int, newPassword string) error
func (r PostgresUserRepository) RequestRole(userId, roleId int) error
func (r PostgresUserRepository) GetRole(userId int) (roleId int, err error)
