package repository

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	SaveUser(username, password, contact string) (int, error)
	GetAllUsers() ([]*domain.User, error)
	GetUserByCredentials(username, password string) (*domain.User, error)
	GetUserById(id int) (*domain.User, error)
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetUserRole(userId int) (*domain.Role, error)
}

type TeamRepository interface {
	SaveTeam(teamLeaderId int, teamName string) (int, error)
	GetTeamById(teamId int) (*domain.Team, error)
	GetTeamByToken(token string) (*domain.Team, error)
	AddMember(teamId, userId int) error
	GetTeams() ([]*domain.Team, error)
}

type TaskRepository interface {
	SaveTask(task domain.Task) (int, error)
	GetTaskById(id int) (*domain.Task, error)
	GetTasks(domain.TaskFilter) ([]*domain.Task, error)
	GetTasksSolvedByTeam(teamId int) ([]*domain.SolvedTask, error)
	UpdateTask(id int, update domain.TaskUpdate) error
	SolveTask(taskId, teamId int) error
	SaveTaskSubmission(submission domain.TaskSubmission) error
	GetTaskSubmissions(taskId, teamId int) ([]*domain.TaskSubmission, error)
	CheckSolvedTask(taskId, teamId int) (bool, error)
	FindOrCreateAuthor(tx *sqlx.Tx, author domain.Author) (int, error)
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
