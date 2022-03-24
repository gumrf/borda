package service

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository, taskRepo repository.TaskRepository) *UserService {
	return &UserService{
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *UserService) IsUserInTeam(userId int) bool {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return false
	}

	if user.Team.Id == 0 {
		return false
	}

	return true
}

func (s *UserService) ShowAllTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *UserService) TryToSolveTask(submission domain.SubmitTaskRequest) (string, error) {
	var task *domain.Task
	var err error

	// Перенести проверку в репозиторий
	//var isTaskSolved bool
	//isTaskSolved, err = s.taskRepo.CheckSolvedTask(submission.TaskId, submission.TeamId)
	//if err != nil {
	//	return "Error on cheking solved task", err
	//}
	//if isTaskSolved {
	//	return "Task already solved!", nil
	//}

	task, err = s.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return "", err
	}

	if submission.Flag == task.Flag {
		err = s.taskRepo.SolveTask(task.Id, submission.TeamId)

		if err != nil {
			return "Error on SolveTask", err
		}

		err = s.taskRepo.SaveTaskSubmission(submission, true)
		if err != nil {
			return "Error on FillTaskSubmission true", err
		}

		return "Submission is correct", nil
	} else {
		err = s.taskRepo.SaveTaskSubmission(submission, false)
		if err != nil {
			return "Error on FillTaskSubmission false", err
		}

		return "Submission is incorrect", nil
	}

}

func (s *UserService) GetTaskSubmissions(input domain.SubmitTaskRequest) ([]*domain.TaskSubmission, error) {
	submissions, err := s.taskRepo.GetTaskSubmissions(input.TaskId, input.UserId)
	if err != nil {
		return nil, err
	}
	return submissions, nil
}
