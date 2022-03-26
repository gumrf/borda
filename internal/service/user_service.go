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

func (s *UserService) ShowAllTasks(filter domain.TaskFilter, userId int) ([]domain.TaskUserResponse, error) {

	// filter.IsActive = true
	// filter.IsDisabled = false ??????????????????????//

	var tasks []*domain.Task
	tasks, err := s.taskRepo.GetTasks(filter) // Получили таски по фильтру
	if err != nil {
		return nil, err
	}

	var teamId int
	teamId, err = s.teamRepo.GetTeamByUserId(userId) //Надо откуда-то волшебным образом высрать TeamId
	if err != nil {
		return nil, err
	}

	userTasks := make([]domain.TaskUserResponse, len(tasks))

	// Ебаный ужас НО ОНО РАБОАЕТ ЭТО ПРОСТО ОХУЕННО
	for i := range tasks {
		userTasks[i].Id = tasks[i].Id
		userTasks[i].Title = tasks[i].Title
		userTasks[i].Description = tasks[i].Description
		userTasks[i].Author = tasks[i].Author
		userTasks[i].Category = tasks[i].Category
		userTasks[i].Hint = tasks[i].Hint
		userTasks[i].Complexity = tasks[i].Complexity
		userTasks[i].Points = tasks[i].Points
		userTasks[i].IsSolved, err = s.taskRepo.CheckSolvedTask(tasks[i].Id, teamId)
		if err != nil {
			return nil, err
		}
		submissions, err := s.taskRepo.GetTaskSubmissions(tasks[i].Id, userId)
		if err != nil {
			return nil, err
		}

		userTasks[i].UserTaskSubmissions = make([]domain.UserTaskSubmission, len(submissions))

		for j := range submissions {
			userTasks[i].UserTaskSubmissions[j].Flag = submissions[j].Flag
			userTasks[i].UserTaskSubmissions[j].Timestemp = submissions[j].Timestemp
			userTasks[i].UserTaskSubmissions[j].UserId = submissions[j].UserId
		}

	}

	return userTasks, nil
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
