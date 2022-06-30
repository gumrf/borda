package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecaseCreateTask struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecaseCreateTask(TaskRepo repository.TaskRepository) *AdminUsecaseCreateTask {
	return &AdminUsecaseCreateTask{
		taskRepo: TaskRepo,
	}
}

func (u *AdminUsecaseCreateTask) Execute(task domain.Task) ([]*domain.Task, error) {
	id, err := u.taskRepo.SaveTask(task)
	if err != nil {
		return nil, domain.ErrTaskCreate
	}

	filter := domain.TaskFilter{
		Id: id,
	}

	createdTask, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTaskNotFound
	}

	return createdTask, nil
}
