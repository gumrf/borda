package service

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserService struct {
	taskRepo repository.TaskRepository
}

func NewUserService(tsr repository.TaskRepository) *UserService {
	return &UserService{taskRepo: tsr}
}

func (u *UserService) ShowAllTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	var tasks []*domain.Task

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *UserService) TryToSolveTask(submission domain.SubmitTaskRequest) (string, error) {
	var task *domain.Task
	var err error

	// Перенести проверку в репозиторий
	//var isTaskSolved bool
	//isTaskSolved, err = a.taskRepo.CheckSolvedTask(submission.TaskId, submission.TeamId)
	//if err != nil {
	//	return "Error on cheking solved task", err
	//}
	//if isTaskSolved {
	//	return "Task already solved!", nil
	//}

	task, err = a.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return "", err
	}

	if submission.Flag == task.Flag {
		err = a.taskRepo.SolveTask(task.Id, submission.TeamId)

		if err != nil {
			return "Error on SolveTask", err
		}

		err = a.taskRepo.SaveTaskSubmission(submission, true)
		if err != nil {
			return "Error on FillTaskSubmission true", err
		}

		return "Submission is correct", nil
	} else {
		err = a.taskRepo.SaveTaskSubmission(submission, false)
		if err != nil {
			return "Error on FillTaskSubmission false", err
		}

		return "Submission is incorrect", nil
	}

}

func (a *UserService) GetTaskSubmissions(input domain.SubmitTaskRequest) ([]*domain.TaskSubmission, error) {
	submissions, err := a.taskRepo.GetTaskSubmissions(input.TaskId, input.UserId)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
