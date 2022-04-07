package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseSubmitFlag struct {
	taskRepo repository.TaskRepository
}

func NewUsecaseSubmitFlag(TaskRepo repository.TaskRepository) *UsecaseSubmitFlag {
	return &UsecaseSubmitFlag{
		taskRepo: TaskRepo,
	}
}

func (u *UsecaseSubmitFlag) Execute(submission domain.TaskSubmission) error {
	task, err := u.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return domain.ErrTaskNotFound
	}

	if submission.Flag == task.Flag {
		submission.IsCorrect = true

		if err := u.taskRepo.SolveTask(task.Id, submission.TeamId); err != nil {
			return domain.ErrTaskSolve
		}
	}

	if err := u.taskRepo.SaveTaskSubmission(submission); err != nil {
		return domain.ErrTaskSaveSubmission
	}

	return nil
}
