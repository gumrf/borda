package team

import (
	"borda/internal/pkg/core"
)

type TeamService struct {
	userRepository core.UserRepository
	teamRepository core.TeamRepository
}

func NewTeamService(userRepository core.UserRepository, teamRepository core.TeamRepository) *TeamService {
	return &TeamService{
		userRepository: userRepository,
		teamRepository: teamRepository,
	}
}

func (ts *TeamService) JoinTeam(userId int, payload string) error {
	// TODO: Validate payload
	// TODO: Wrap in transaction
	_, err := ts.teamRepository.FindByToken(payload)
	if err != nil {
		return err
	}

	// if err := ts.teamRepository.AddMember(team.Id, userId); err != nil {
	// 	return err
	// }

	return nil
}

func (ts *TeamService) CreateTeam(userId int, payload string) error {
	// TODO: Validate payload
	// TODO: Wrap in transaction

	// uuid := uuid.New().String()

	// if _, err := us.teamRepository.Save(); err != nil {
	// 	return err
	// }
	// if err := ts.teamRepository.AddMember(team.Id, userId); err != nil {
	// 	return err
	// }
	return nil
}
