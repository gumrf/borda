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

func (a *AdminService) GetAllTasks() ([]*domain.Task, error) {
	var filter domain.TaskFilter

	tasks, err := a.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *AdminService) CreateNewTask(task domain.Task) ([]*domain.Task, error) {

	id, err := a.taskRepo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	filter := domain.TaskFilter{
		Id: id,
	}

	createdTask, err := a.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (a *AdminService) UpdateTask(taskId int, taskUpdate domain.TaskUpdate) error {
	return a.taskRepo.UpdateTask(taskId, taskUpdate)
}
