package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecaseGetAllTasks struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecaseGetAlTasks(TaskRepo repository.TaskRepository) *AdminUsecaseGetAllTasks {
	return &AdminUsecaseGetAllTasks{
		taskRepo: TaskRepo,
	}
}

func (u *AdminUsecaseGetAllTasks) Execute() ([]*domain.Task, error) {
	var filter domain.TaskFilter

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTasksNotFound
	}

	return tasks, nil
}
