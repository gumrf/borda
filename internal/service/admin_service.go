package service

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type AdminService struct {
	taskRepo repository.TaskRepository
}

func NewAdminService(tsr repository.TaskRepository) *AdminService {
	return &AdminService{taskRepo: tsr}
}

func (a *AdminService) CreateNewTask(task domain.Task) (int, error) {

	id, err := a.taskRepo.CreateNewTask(task)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (a *AdminService) UpdateTask(taskId int, dataForUpdate domain.TaskUpdate) error {
	return a.taskRepo.UpdateTask(taskId, dataForUpdate)
}
