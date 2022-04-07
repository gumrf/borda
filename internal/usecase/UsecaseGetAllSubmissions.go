package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseGetAllSubmissions struct {
	taskRepo repository.TaskRepository
}

func NewUsecaseGetAllSubmissions(TaskRepo repository.TaskRepository) *UsecaseGetAllSubmissions {
	return &UsecaseGetAllSubmissions{
		taskRepo: TaskRepo,
	}
}

func (u *UsecaseGetAllSubmissions) Execute(taskId, userId int) ([]*domain.TaskSubmission, error) {
	submissions, err := u.taskRepo.GetTaskSubmissions(taskId, userId)
	if err != nil {
		return nil, domain.ErrTaskSubmissionsNotFound
	}
	return submissions, nil
}
