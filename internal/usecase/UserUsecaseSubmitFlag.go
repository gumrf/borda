package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecaseSubmitFlag struct {
	taskRepo repository.TaskRepository
}

func NewUserUsecaseSubmitFlag(TaskRepo repository.TaskRepository) *UserUsecaseSubmitFlag {
	return &UserUsecaseSubmitFlag{
		taskRepo: TaskRepo,
	}
}

func (u *UserUsecaseSubmitFlag) Execute(submission domain.TaskSubmission) (domain.SubmitFlagResponse, error) {
	var response domain.SubmitFlagResponse

	task, err := u.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return response, domain.ErrTaskNotFound
	}

	if submission.Flag == task.Flag {
		submission.IsCorrect = true

		if err := u.taskRepo.SolveTask(task.Id, submission.TeamId); err != nil {
			return response, domain.ErrTaskSolve
		}
	}

	if err := u.taskRepo.SaveTaskSubmission(submission); err != nil {
		return response, domain.ErrTaskSaveSubmission
	}

	response.TaskId = task.Id
	response.IsCorrect = submission.IsCorrect

	return response, nil
}
