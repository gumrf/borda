package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecaseGetProfile struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserUsecaseGetProfile(UserRepo repository.UserRepository,
	TeamRepo repository.TeamRepository) *UserUsecaseGetProfile {
	return &UserUsecaseGetProfile{
		userRepo: UserRepo,
		teamRepo: TeamRepo,
	}
}

func (u *UserUsecaseGetProfile) Execute(userId int) (domain.PrivateUserProfileResponse, error) {
	user, err := u.userRepo.GetUserById(userId)
	if err != nil {
		return domain.PrivateUserProfileResponse{}, err
	}

	team, err := u.teamRepo.GetTeamById(user.TeamId)
	if err != nil {
		return domain.PrivateUserProfileResponse{}, err
	}

	response := domain.PrivateUserProfileResponse{
		Id:       user.Id,
		Username: user.Username,
		Contact:  user.Contact,
		Team: domain.TeamResponse{
			Id:    team.Id,
			Name:  team.Name,
			Token: team.Token,
			Captain: func() domain.TeamMember {
				var captain domain.TeamMember
				for _, m := range team.Members {
					if m.UserId == team.TeamLeaderId {
						captain.UserId = m.UserId
						captain.Name = m.Name
					}
				}
				return captain
			}(),
			Members: team.Members,
		},
	}

	return response, nil
}
