package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecaseGetScoreboard struct {
	teamRepo repository.TeamRepository
	taskRepo repository.TaskRepository
}

func NewUserUsecaseGetScoreboard(TeamRepo repository.TeamRepository,
	TaskRepo repository.TaskRepository) *UserUsecaseGetScoreboard {
	return &UserUsecaseGetScoreboard{
		teamRepo: TeamRepo,
		taskRepo: TaskRepo,
	}
}

func (uc *UserUsecaseGetScoreboard) Execute() ([]domain.Scoreboard, error) {
	teams, err := uc.teamRepo.GetTeams()
	if err != nil {
		return nil, err
	}

	response := make([]domain.Scoreboard, 0)

	for _, team := range teams {
		solvedTasks, err := uc.taskRepo.GetTasksSolvedByTeam(team.Id)
		if err != nil {
			return nil, err
		}

		var score int
		for _, task := range solvedTasks {
			task, err := uc.taskRepo.GetTaskById(task.TaskId)
			if err != nil {
				return nil, err
			}

			if !task.IsDisabled {
				score += task.Points
			}
		}

		responseItem := domain.Scoreboard{
			TeamName:         team.Name,
			Score:            score,
			TeamMembersCount: len(team.Members),
		}

		response = append(response, responseItem)
	}

	return response, nil
}
