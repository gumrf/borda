package core

import (
	"borda/internal/core/entity"
)

// UserRepository
type UserRepository interface {
	Create(username, password, contact string) (int, error)
	// Find(username, password string) (*entity.User, error)
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetRole(userId int) (entity.Role, error)
}

// TemaRepository
type TeamRepository interface {
	Create(teamLeaderId int, teamName string) (entity.Team, error)
	Get(teamId int) (team entity.Team, err error)
	AddMember(teamId, userId int) error
	GetMembers(teamId int) (users []entity.User, err error)
}

// TaskRepository
type TaskRepository interface {
	CreateNewTask(task entity.Task) (int, error)
	GetTaskById(id int) (*entity.Task, error)
	GetTasks(entity.TaskFilter) ([]*entity.Task, error)
	UpdateTask(id int, update entity.TaskUpdate) error
	SolveTask(taskId, teamId int) error
}

// SettingsRepository
type SettingsRepository interface {
	Get(key string) (value string, err error)
	Set(key string, value string) (settingId int, err error)
}

// Repository
type Repository interface {
	Users() UserRepository
	Settings() SettingsRepository
	Teams() TeamRepository
	Tasks() TaskRepository
}
