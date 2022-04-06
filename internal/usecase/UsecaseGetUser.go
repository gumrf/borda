package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UsecaseGetUser struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUsecaseGetUser(UserRepo repository.UserRepository,
	TeamRepo repository.TeamRepository) *UsecaseGetUser {
	return &UsecaseGetUser{
		userRepo: UserRepo,
		teamRepo: TeamRepo,
	}
}

func (u *UsecaseGetUser) Execute(userId int, me bool) (domain.UserProfileResponse, error) {
	var result domain.UserProfileResponse
	switch me {
	case true:
		user, err := u.userRepo.GetUserById(userId)
		if err != nil {
			return domain.UserProfileResponse{}, err
		}

		team, err := u.teamRepo.GetTeamById(user.TeamId)
		if err != nil {
			return domain.UserProfileResponse{}, domain.ErrTeamNotFound
		}

		teamMembers, err := u.teamRepo.GetMembers(user.TeamId)
		if err != nil {
			return domain.UserProfileResponse{}, domain.ErrMembers
		}

		membersResponse := make([]domain.MemberResponse, 0)

		for _, teamMember := range teamMembers {
			isCaptain := false
			if teamMember.Id == team.TeamLeaderId {
				isCaptain = true
			}

			memberResponse := domain.MemberResponse{
				Username:  teamMember.Username,
				IsCaptain: isCaptain,
			}

			membersResponse = append(membersResponse, memberResponse)
		}

		userProfileResponse := domain.UserProfileResponse{
			Id:          userId,
			Username:    user.Username,
			Contact:     user.Contact,
			TeamName:    team.Name,
			TeamId:      team.Id,
			TeamMembers: membersResponse,
		}

		result, err = userProfileResponse, nil
	default:
		user, err := u.userRepo.GetUserById(userId)
		if err != nil {
			return domain.UserProfileResponse{}, domain.ErrUserNotFound
		}

		team, err := u.teamRepo.GetTeamById(user.TeamId)
		if err != nil {
			return domain.UserProfileResponse{}, domain.ErrTeamNotFound
		}

		teamMembers, err := u.teamRepo.GetMembers(user.TeamId)
		if err != nil {
			return domain.UserProfileResponse{}, domain.ErrMembers
		}

		membersResponse := make([]domain.MemberResponse, 0)

		for _, teamMember := range teamMembers {
			memberResponse := domain.MemberResponse{
				Username: teamMember.Username,
			}
			membersResponse = append(membersResponse, memberResponse)
		}

		userProfileResponse := domain.UserProfileResponse{
			Id:          user.Id,
			Username:    user.Username,
			TeamId:      team.Id,
			TeamName:    team.Name,
			TeamMembers: membersResponse,
		}

		result, err = userProfileResponse, nil
		return userProfileResponse, nil
	}

	return result, nil
}
