package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecase struct {
	taskRepo repository.TaskRepository
}

func NewUserUsecase(tsr repository.TaskRepository) *UserUsecase {
	return &UserUsecase{taskRepo: tsr}
}

func (u *UserUsecase) ShowAllTasks() ([]*domain.Task, error) {
	var tasks []*domain.Task
	var filter domain.TaskFilter

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
