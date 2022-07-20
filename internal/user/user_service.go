package user

import (
	"borda/internal/pkg/core"
	"borda/pkg/transaction"
	"context"
	"errors"
)

type UserService struct {
	userRepository core.UserRepository
	teamRepository core.TeamRepository
	txManager      transaction.Manager
}

func NewUserService(userRepository core.UserRepository,
	teamRepository core.TeamRepository, txManager transaction.Manager) *UserService {

	return &UserService{
		userRepository: userRepository,
		teamRepository: teamRepository,
		txManager:      txManager,
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
	var actualUser User

	err := us.txManager.Run(
		context.Background(),
		func(ctx context.Context) error {
			user, err := us.userRepository.FindById(ctx, id)
			actualUser = user

			if err != nil {
				return err
			}

			return nil
		},
		true,
	)

	if err != nil {
		return UserResponse{}, err
	}

	// team, err := us.teamRepository.FindById(user.TeamId)

	response := UserResponse{
		Id:       actualUser.Id,
		Username: actualUser.Username,
	}

	return response, nil
}

func (us *UserService) GetAllUsers() ([]UserResponse, error) {
	var actualUsers []User

	err := us.txManager.Run(
		context.Background(),
		func(ctx context.Context) error {
			users, err := us.userRepository.FindAll(ctx)
			actualUsers = users
			if err != nil {
				// Wrap error!!!!
				return err
			}

			return nil
		},
		true,
	)

	if err != nil {
		// Wrap error!!!!
		return []UserResponse{}, err
	}

	response := make([]UserResponse, 0)

	for _, user := range actualUsers {

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
