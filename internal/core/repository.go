package core

import (
	"borda/internal/core/entity"
)

// UserRepository represent op
type UserRepository interface {
	Create(username, password, contact string) (userId int, err error)
	UpdatePassword(userId int, newPassword string) error
	AssignRole(userId, roleId int) error
	GetRole(userId int) (entity.Role, error)
}

// TemaRepository
type TeamRepository interface {
	Create(teamLeaderId int, teamName string) (team entity.Team, err error)
	Get(teamId int) (team entity.Team, err error)
	AddMember(teamId, userId int) error
	GetMembers(teamId int) (users []entity.User, err error) // TODO implement
}

// Taskrepository
type TaskRepository interface {
	Get(taskId int) (entity.Task, error)
	FindTasks(entity.TaskFilter) ([]entity.Task, error)
	Solve(taskId, teamId int) error
	Create(task entity.Task) (taskId int, err error)
	Update(taskId int, newTask entity.Task) error
}

type SettingsRepository interface {
	Get(key string) (value string, err error)
	Set(key string, value string) (settingsId int, err error)
}

type Repository interface {
	Users() UserRepository
	Settings() SettingsRepository
	Teams() TeamRepository
	Tasks() TaskRepository
}
