package service

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(userRepo repository.UserRepository, taskRepo repository.TaskRepository,
	teamRepo repository.TeamRepository) *UserService {
	return &UserService{
		taskRepo: taskRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *UserService) GetUser(id int) (*domain.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *UserService) IsUserInTeam(userId int) (int, bool) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return 0, false
	}

	if user.TeamId <= 0 {
		return 0, false
	}

	return user.TeamId, true
}

func (s *UserService) GetAllUsers() ([]domain.UserResponse, error) {
	//Получил всех пользовотелей, которые в командах
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
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
			team, err := s.teamRepo.GetTeamById(user.TeamId)
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
