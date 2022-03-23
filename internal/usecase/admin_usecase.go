package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminUsecase struct {
	taskRepo repository.TaskRepository
}

func NewAdminUsecase(tsr repository.TaskRepository) *AdminUsecase {
	return &AdminUsecase{taskRepo: tsr}
}

func (a *AdminUsecase) CreateNewTask(task domain.Task) (int, error) {

	id, err := a.taskRepo.CreateNewTask(task)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (a *AdminUsecase) UpdateTask(taskId int, dataForUpdate domain.TaskUpdate) error {
	return a.taskRepo.UpdateTask(taskId, dataForUpdate)
}
