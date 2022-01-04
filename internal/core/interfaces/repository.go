package interfaces

import "borda/internal/core/entities"

type TaskRepository interface {
	// CreateNewTask создает новый такс
	CreateNewTask(t entities.Task) (entities.Task, error)
	// DeleteTask удаляет таск
	DeleteTask(taskId int) error
	// Update обновляет
	UpdateTask(t entities.Task) (entities.Task, error)
	// Show открывает все скрытые таски
	Show()
	// Открывает определенный таск
	ShowOne(taskId int) error
	// Close скрывает все таски
	Close()
	// CloseOne скрывает определенный таск
	CloseOne(taskId int) error
	// Disable убирает таск
	Disable(taskId int) error
	// Enable возвращает таск
	Enable(taskId int) error
	// Backup выгружает все таски в json формате
	Backup() (filename string, err error)
	// Import загружает таски из json или markdown файла
	Import(file string) (taskCount int, err error)
	// Getauthors return Authors associeted with task
	GetAuthors(taskId int) ([]entities.Author, error)
}

type UserRepository interface{}

type TeamRepository interface{}

type RoleRepository interface{}

type Repository interface {
	Users() (UserRepository, error)
	Roles() (RoleRepository, error)
	Tasks() (TaskRepository, error)
	Teams() (TeamRepository, error)
}
