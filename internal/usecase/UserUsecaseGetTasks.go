package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
	"fmt"
)

type UserUsecaseGetTasks struct {
	taskRepo repository.TaskRepository
	teamRepo repository.TeamRepository
}

func NewUserUsecaseGetTasks(TaskRepo repository.TaskRepository,
	TeamRepo repository.TeamRepository) *UserUsecaseGetTasks {
	return &UserUsecaseGetTasks{
		taskRepo: TaskRepo,
		teamRepo: TeamRepo,
	}
}

func (uc *UserUsecaseGetTasks) Execute(teamId int) ([]domain.PublicTaskResponse, error) {
	filter := domain.TaskFilter{
		IsActive:   true,
		IsDisabled: false,
	}

	tasks, err := uc.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, domain.ErrTasksNotFound
	}

	team, err := uc.teamRepo.GetTeamById(teamId)
	if err != nil {
		return nil, fmt.Errorf("GetTeamById: %v", err)
	}

	findUsernameById := func(id int) string {
		for _, member := range team.Members {
			if member.UserId == id {
				return member.Name
			}
		}
		return ""
	}

	response := make([]domain.PublicTaskResponse, 0)

	for _, task := range tasks {
		submissions, err := uc.taskRepo.GetTaskSubmissions(task.Id, team.Id)
		if err != nil {
			return nil, domain.ErrTaskSubmissionsNotFound
		}

		isTaskSolved, err := uc.taskRepo.CheckSolvedTask(task.Id, team.Id)
		if err != nil {
			return nil, domain.ErrTaskSubmissionsNotFound
		}

		submissionsResponse := make([]domain.SubmissionResponse, 0)
		for _, submission := range submissions {
			submissionResponseItem := domain.SubmissionResponse{
				Username:  findUsernameById(submission.UserId),
				Flag:      submission.Flag,
				IsCorrect: submission.IsCorrect,
				Timestamp: submission.Timestamp,
			}

			submissionsResponse = append(submissionsResponse, submissionResponseItem)
		}

		responseItem := domain.PublicTaskResponse{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Link:        task.Link,
			Category:    task.Description,
			Complexity:  task.Category,
			Points:      task.Points,
			Hint:        task.Hint,
			IsSolved:    isTaskSolved,
			Submissions: submissionsResponse,
			Author:      task.Author,
		}

		response = append(response, responseItem)
	}

	return response, nil
}
