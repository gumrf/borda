package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecaseCreateNewTask struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecaseCreateNewTask(TaskRepo repository.TaskRepository) *AdminUsecaseCreateNewTask {
	return &AdminUsecaseCreateNewTask{
		taskRepo: TaskRepo,
	}
}

func (u *AdminUsecaseCreateNewTask) Execute(task domain.Task) ([]*domain.Task, error) {
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
