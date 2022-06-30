package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecaseUpdateTask struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecaseUpdateTask(TaskRepo repository.TaskRepository) *AdminUsecaseUpdateTask {
	return &AdminUsecaseUpdateTask{
		taskRepo: TaskRepo,
	}
}

func (u *AdminUsecaseUpdateTask) Execute(taskId int, taskUpdate domain.TaskUpdate) ([]*domain.Task, error) {
	if err := u.taskRepo.UpdateTask(taskId, taskUpdate); err != nil {
		return nil, domain.ErrTaskUpdate
	}

	filter := domain.TaskFilter{
		Id: taskId,
	}

	updatedTask, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTaskNotFound
	}

	return updatedTask, nil
}
