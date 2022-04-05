package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseGetAllUsers struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUsecaseGetAllUsers(UserRepo repository.UserRepository,
	TeamRepo repository.TeamRepository) *UsecaseGetAllUsers {
	return &UsecaseGetAllUsers{
		userRepo: UserRepo,
		teamRepo: TeamRepo,
	}
}

func (u *UsecaseGetAllUsers) Execute() ([]domain.UserResponse, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, domain.ErrUsersNotFound
	}

	usersResponse := make([]domain.UserResponse, 0)

	//Прохожусь по всем юзерам отдельно
	for _, user := range users {
		var userResponse domain.UserResponse

		if user.TeamId == 0 {
			userResponse = domain.UserResponse{
				Username: user.Username,
			}

		} else {
			team, err := u.teamRepo.GetTeamById(user.TeamId)
			if err != nil {
				return nil, err
			}

			userResponse = domain.UserResponse{
				Username: user.Username,
				TeamName: team.Name,
			}
		}
		usersResponse = append(usersResponse, userResponse)
	}

	return usersResponse, nil
}
