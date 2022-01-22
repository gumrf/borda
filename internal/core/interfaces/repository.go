package interfaces

import "borda/internal/core/entity"

type UserRepository interface {
	Create(username, password, contact string) (userId int, err error)
	UpdatePassword(userId int, newPassword string) error
	RequestRole(userId, roleId int) error
	GetRole(userId int) (roleId int, err error)
}
type TeamRepository interface {
	Create(teamLeaderId int, teamName string) (team entity.Team, err error)
	AddMember(teamId, userId int) error
	Get(teamId int) (team entity.Team, err error)
}

type TaskRepository interface {
	Get(taskId int) (entity.Task, error)
	GetMany(taskParams interface{}) ([]entity.Task, error)
	Solve(taskId int) error
	Save(task entity.Task) (taskId int, err error)
	Update(oldTask, newTask entity.Task) error
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
