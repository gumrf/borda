package service

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(userRepo repository.UserRepository, taskRepo repository.TaskRepository,
	teamRepo repository.TeamRepository) *UserService {
	return &UserService{
		taskRepo: taskRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
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

func (s *UserService) GetAllTasks(userId int) ([]domain.TaskUserResponse, error) {
	filter := domain.TaskFilter{
		IsActive:   true,
		IsDisabled: false,
	}

	var tasks []*domain.Task
	// Получили таски по фильтру
	tasks, err := s.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	var teamId int
	//Надо откуда-то волшебным образом высрать TeamId
	teamId, err = s.teamRepo.GetTeamByUserId(userId)
	if err != nil {
		return nil, err
	}

	//Получаю всех участников команды
	users, err := s.teamRepo.GetMembers(teamId)
	if err != nil {
		return nil, err
	}

	//Вношу имена в мапу [user.id]username
	usernames := make(map[int]string)
	for _, user := range users {
		usernames[user.Id] = user.Username
	}

	userScopedTasks := make([]domain.TaskUserResponse, 0)

	for _, task := range tasks {

		//Получаю все решения этого пользователя
		allSubmissions, err := s.taskRepo.GetTaskSubmissions(task.Id, teamId)
		if err != nil {
			return nil, err
		}

		//Привожу решения пользователя в вид для этого пользователя
		taskSubmissionResponse := make([]domain.TaskSubmissionResponse, 0)
		for _, sub := range allSubmissions {
			submissionResponse := domain.TaskSubmissionResponse{
				Username:  usernames[sub.UserId],
				Flag:      sub.Flag,
				IsCorrect: sub.IsCorrect,
				Timestemp: sub.Timestemp,
			}
			taskSubmissionResponse = append(taskSubmissionResponse, submissionResponse)
		}

		//Проверка, решен ли таск
		IsSolved, err := s.taskRepo.CheckSolvedTask(task.Id, teamId)
		if err != nil {
			return nil, err
		}

		//Заполнение формы таска для юзера
		taskResponse := domain.TaskUserResponse{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Category:    task.Description,
			Complexity:  task.Category,
			Points:      task.Points,
			Hint:        task.Hint,
			IsSolved:    IsSolved,
			Submissions: taskSubmissionResponse,
			Author:      task.Author,
		}

		userScopedTasks = append(userScopedTasks, taskResponse)

	}

	return userScopedTasks, nil
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
