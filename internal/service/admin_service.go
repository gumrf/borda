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

func (s *AdminService) GetAllTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *AdminService) CreateNewTask(task domain.Task) (int, error) {

	id, err := a.taskRepo.SaveTask(task)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (a *AdminService) UpdateTask(taskId int, taskUpdate domain.TaskUpdate) error {
	return a.taskRepo.UpdateTask(taskId, taskUpdate)
}
