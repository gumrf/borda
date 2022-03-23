package repository

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateNewUser(username, password, contact string) (int, error)
	FindUserByCredentials(username, password string) (*domain.User, error)
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetRole(userId int) (domain.Role, error)
	IsUsernameExists(username string) error
	// SetSession(userId int, session domain.Session) error
}

type TeamRepository interface {
	CreateNewTeam(teamLeaderId int, teamName string) (int, error)
	GetTeamById(teamId int) (domain.Team, error)
	GetTeamByToken(token string) (domain.Team, error)
	AddMember(teamId, userId int) error
	GetMembers(teamId int) ([]domain.User, error)
	IsTeamNameExists(teamName string) error
	IsTeamTokenValid(token string) error
	IsTeamFull(teamId int) error
}

type TaskRepository interface {
	CreateNewTask(task domain.Task) (int, error)
	GetTaskById(id int) (*domain.Task, error)
	GetTasks(domain.TaskFilter) ([]*domain.Task, error)
	UpdateTask(id int, update domain.TaskUpdate) error
	SolveTask(taskId, teamId int) error
	FillTaskSubmission(value domain.SubmitTaskRequest, isCorrect bool) error
	ShowTaskSubmissions(value domain.SubmitTaskRequest) ([]*domain.TaskSubmission, error)
	ChekSolvedTask(taskId, teamId int) (bool, error)
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
