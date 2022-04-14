package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecaseGetTasks struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecaseGetTasks(TaskRepo repository.TaskRepository) *AdminUsecaseGetTasks {
	return &AdminUsecaseGetTasks{
		taskRepo: TaskRepo,
	}
}

func (u *AdminUsecaseGetTasks) Execute() ([]*domain.Task, error) {
	var filter domain.TaskFilter

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTasksNotFound
	}

	return tasks, nil
}
