package repository

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(username, password, contact string) (int, error)
	FindUser(username, password string) (*domain.User, error)
	FindUserByUsename(username string) error
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetRole(userId int) (domain.Role, error)
}

type TeamRepository interface {
	Create(teamLeaderId int, teamName string) (domain.Team, error)
	Get(teamId int) (team domain.Team, err error)
	AddMember(teamId, userId int) error
	GetMembers(teamId int) (users []domain.User, err error)
}

type TaskRepository interface {
	CreateNewTask(task domain.Task) (int, error)
	GetTaskById(id int) (*domain.Task, error)
	GetTasks(domain.TaskFilter) ([]*domain.Task, error)
	UpdateTask(id int, update domain.TaskUpdate) error
	SolveTask(taskId, teamId int) error
}

type SettingsRepository interface {
	Get(key string) (value string, err error)
	Set(key string, value string) (settingId int, err error)
}

type Repository struct {
	db       *sqlx.DB
	Users    UserRepository
	Teams    TeamRepository
	Tasks    TaskRepository
	Settings SettingsRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:       db,
		Users:    postgres.NewUserRepository(db),
		Teams:    postgres.NewTeamRepository(db),
		Tasks:    postgres.NewTaskRepository(db),
		Settings: postgres.NewSettingsRepository(db),
	}
}
