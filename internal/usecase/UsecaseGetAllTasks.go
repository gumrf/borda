package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseGetAllTasks struct {
	taskRepo repository.TaskRepository
	teamRepo repository.TeamRepository
}

func NewUsecaseGetAllTasks(TaskRepo repository.TaskRepository,
	TeamRepo repository.TeamRepository) *UsecaseGetAllTasks {
	return &UsecaseGetAllTasks{
		taskRepo: TaskRepo,
		teamRepo: TeamRepo,
	}
}

func (u *UsecaseGetAllTasks) Execute(id int) ([]domain.UserTaskResponse, error) {
	filter := domain.TaskFilter{
		IsActive:   true,
		IsDisabled: false,
	}

	// Получили таски по фильтру
	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTasksNotFound
	}

	//Получаю всех участников команды
	users, err := u.teamRepo.GetMembers(id)
	if err != nil {
		return nil, domain.ErrTeamMembersNotFound
	}

	//Вношу имена в мапу [user.id]username
	usernames := make(map[int]string, 4)
	for _, user := range users {
		usernames[user.Id] = user.Username
	}

	userTasksResponse := make([]domain.UserTaskResponse, 0)

	for _, task := range tasks {

		// Get team submissions for task
		submissions, err := u.taskRepo.GetTaskSubmissions(task.Id, id)
		if err != nil {
			return nil, domain.ErrTaskSubmissionsNotFound
		}

		// Build submissions response
		taskSubmissionsResponse := make([]domain.TaskSubmissionResponse, 0)
		for _, submission := range submissions {
			// Allocate sibmission object
			submissionResponse := domain.TaskSubmissionResponse{
				Username:  usernames[submission.UserId],
				Flag:      submission.Flag,
				IsCorrect: submission.IsCorrect,
				Timestamp: submission.Timestamp,
			}
			// Append allocated submission object to array of submissions
			taskSubmissionsResponse = append(taskSubmissionsResponse, submissionResponse)
		}

		// Check if task solved
		isTaskSolved, err := u.taskRepo.CheckSolvedTask(task.Id, id)
		if err != nil {
			return nil, domain.ErrTaskSubmissionsNotFound
		}

		// Allocate and fill task object
		taskResponse := domain.UserTaskResponse{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Category:    task.Description,
			Complexity:  task.Category,
			Points:      task.Points,
			Hint:        task.Hint,
			IsSolved:    isTaskSolved,
			Submissions: taskSubmissionsResponse,
			Author:      task.Author,
		}

		// Append task object to UserTasksResponse array
		userTasksResponse = append(userTasksResponse, taskResponse)
	}

	return userTasksResponse, nil
}
