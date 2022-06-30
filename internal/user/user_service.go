package user

import (
	"borda/internal/pkg/core"
	"errors"
)

type UserService struct {
	userRepository core.UserRepository
	teamRepository core.TeamRepository
}

func NewUserService(userRepository core.UserRepository,
	teamRepository core.TeamRepository) *UserService {

	return &UserService{
		userRepository: userRepository,
		teamRepository: teamRepository,
	}
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Team     struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"team"`
}

func (us *UserService) GetUser(id int) (UserResponse, error) {
	user, err := us.userRepository.FindById(id)
	if err != nil {
		return UserResponse{}, err
	}

	// team, err := us.teamRepository.FindById(user.TeamId)

	response := UserResponse{
		Id:       user.Id,
		Username: user.Username,
	}

	return response, nil
}

func (us *UserService) GetAllUsers() ([]UserResponse, error) {
	users, err := us.userRepository.FindAll()
	if err != nil {
		// Wrap error!!!!
		return []UserResponse{}, err
	}

	response := make([]UserResponse, 0)

	for _, user := range users {

		// TODO
		// team, err := us.teamRepository.FindById(user.TeamId)

		response = append(response, UserResponse{
			Id:       user.Id,
			Username: user.Username,
		})
	}

	return response, nil
}

type TeamResponse struct{}

type UserProfileResponse struct {
	Id       int          `json:"id"`
	Username string       `json:"username"`
	Contact  string       `json:"contact"`
	Team     TeamResponse `json:"team"`
}

func (us *UserService) GetUserProfile(id int) (UserProfileResponse, error) {
	return UserProfileResponse{}, errors.New("not implemented")
}
