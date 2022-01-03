package repository

import "borda/internal/domain"

type TaskRepositoryI interface {
	// CreateNewTask создает новый такс
	CreateNewTask(t domain.Task) (domain.Task, error)
	// DeleteTask удаляет таск
	DeleteTask(taskId int) error
	// Update обновляет
	UpdateTask(t domain.Task) (domain.Task, error)
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
	GetAuthors(taskId int)([]domain.Author, error)
}

type UserRepositoryI interface {}

type TeamRepositoryI interface{}

type RoleRepositoryI interface{}

type RepositoryI interface {
	Users() (UserRepositoryI, error)
	Roles() (RoleRepositoryI, error)
	Tasks() (TaskRepositoryI, error)
	Teams() (TeamRepositoryI, error)
}
