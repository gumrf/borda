package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseGetScoreboard struct {
	teamRepo repository.TeamRepository
	taskRepo repository.TaskRepository
}

func NewUsecaseScoreboard(TeamRepo repository.TeamRepository,
	TaskRepo repository.TaskRepository) *UsecaseGetScoreboard {
	return &UsecaseGetScoreboard{
		teamRepo: TeamRepo,
		taskRepo: TaskRepo,
	}
}

func (uc *UsecaseGetScoreboard) Execute() ([]domain.Scoreboard, error) {

	//Get all teams with members
	teams, err := uc.teamRepo.GetTeams()
	if err != nil {
		return nil, domain.ErrTeamsNotFound
	}

	response := make([]domain.Scoreboard, 0)

	for _, team := range teams {
		// get all solved tasks for team
		solvedTasks, err := uc.taskRepo.GetSolvedTasks(team.Id)
		if err != nil {
			return nil, domain.ErrSolvedTasksNotFound
		}

		var score int

		// Count score
		for _, solvedTask := range solvedTasks {
			task, err := uc.taskRepo.GetTaskById(solvedTask.TaskId)
			if err != nil {
				return nil, err
			}

			// If task enabled ++
			if !task.IsDisabled {
				score += task.Points
			}
		}

		responseItem := domain.Scoreboard{
			Name:         team.Name,
			Score:        score,
			CountMembers: len(team.Members),
		}

		response = append(response, responseItem)
	}

	return response, nil
}
