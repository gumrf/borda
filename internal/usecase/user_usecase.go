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

func (u *UserUsecase) ShowAllTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	var tasks []*domain.Task

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *UserUsecase) TryToSolveTask(submission domain.SubmitTaskRequest) (string, error) {
	var task *domain.Task
	var err error

	task, err = a.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return "", err
	}

	if submission.Flag == task.Flag {
		err = a.taskRepo.SolveTask(task.Id, submission.TeamId)

		if err != nil {
			return "Error on SolveTask", err
		}

		err = a.taskRepo.FillTaskSubmission(submission, true)
		if err != nil {
			return "Error on FillTaskSubmission true", err
		}

		return "Submission is correct", nil
	} else {
		err = a.taskRepo.FillTaskSubmission(submission, false)
		if err != nil {
			return "Error on FillTaskSubmission false", err
		}

		return "Submission is incorrect", nil
	}

}

func (a *UserUsecase) ShowAllSubmisiions(input domain.SubmitTaskRequest) ([]*domain.TaskSubmission, error) {
	submissions, err := a.taskRepo.ShowTaskSubmissions(input)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
