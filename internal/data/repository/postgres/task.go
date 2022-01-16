package repository

import (
	"borda/internal/core/entity"
	"borda/internal/core/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostgresTaskRepository struct {
	db *sqlx.DB
}

var _ interfaces.TaskRepository = (*PostgresTaskRepository)(nil)

func NewPostgresTaskRepository(db *sqlx.DB) interfaces.TaskRepository {
	return PostgresTaskRepository{db: db}
}

func (r PostgresTaskRepository) Get(taskId int) (entity.Task, error) {
	return entity.Task{}, nil
}
func (r PostgresTaskRepository) GetMany(taskParams interface{}) ([]entity.Task, error) {
	return []entity.Task{}, nil
}
func (r PostgresTaskRepository) Solve(taskId int) error {
	return nil
}
func (r PostgresTaskRepository) Save(task entity.Task) (taskId int, err error) {
	return 0, nil
}
func (r PostgresTaskRepository) Update(oldTask, newTask entity.Task) error {
	return nil
}
