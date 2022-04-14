package usecase

import (
	"borda/internal/repository"
)

type UserUsecaseJoinTeam struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserUsecaseJoinTeam(UserRepo repository.UserRepository,
	TeamRepo repository.TeamRepository) *UserUsecaseJoinTeam {
	return &UserUsecaseJoinTeam{
		userRepo: UserRepo,
		teamRepo: TeamRepo,
	}
}

func (u *UserUsecaseJoinTeam) Execute(userId int, method string, attribute string) error {
	switch method {
	case "create":
		if _, err := u.teamRepo.SaveTeam(userId, attribute); err != nil {
			return err
		}
	case "join":
		team, err := u.teamRepo.GetTeamByToken(attribute)
		if err != nil {
			return err
		}
	
		if err := u.teamRepo.AddMember(team.Id, userId); err != nil {
			return err
		}
	}
	return nil
}