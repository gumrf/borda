package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecaseGetUsers struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserUsecaseGetUsers(UserRepo repository.UserRepository,
	TeamRepo repository.TeamRepository) *UserUsecaseGetUsers {
	return &UserUsecaseGetUsers{
		userRepo: UserRepo,
		teamRepo: TeamRepo,
	}
}

func (uc *UserUsecaseGetUsers) Execute(id ...int) ([]domain.PublicUserProfileResponse, error) {
	users := make([]*domain.User, 0)

	if len(id) > 0 {
		user, err := uc.userRepo.GetUserById(id[0])
		if err != nil {
			return []domain.PublicUserProfileResponse{}, err
		}
		users = append(users, user)
	} else {
		_users, err := uc.userRepo.GetAllUsers()
		if err != nil {
			return nil, domain.ErrUsersNotFound
		}

		users = _users
	}

	response := make([]domain.PublicUserProfileResponse, 0)

	for _, user := range users {
		var responseItem domain.PublicUserProfileResponse

		if user.TeamId == 0 {
			responseItem = domain.PublicUserProfileResponse{
				Id:       user.Id,
				Username: user.Username,
			}

		} else {
			team, err := uc.teamRepo.GetTeamById(user.TeamId)
			if err != nil {
				return nil, err
			}

			responseItem = domain.PublicUserProfileResponse{
				Id:       user.Id,
				Username: user.Username,
				Team: struct {
					Id   int    `json:"id,omitempty"`
					Name string `json:"name,omitempty"`
				}{team.Id, team.Name},
			}
		}
		response = append(response, responseItem)
	}

	return response, nil
}
