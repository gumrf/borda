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

func (s *UserService) GetAllTasks(userId int) ([]domain.UserTaskResponse, error) {
	filter := domain.TaskFilter{
		IsActive:   true,
		IsDisabled: false,
	}

	// Получили таски по фильтру
	tasks, err := s.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	// TODO: Write teamid  to context when che
	//Надо откуда-то волшебным образом высрать TeamId
	teamId, err := s.teamRepo.GetTeamByUserId(userId)
	if err != nil {
		return nil, err
	}

	//Получаю всех участников команды
	users, err := s.teamRepo.GetMembers(teamId)
	if err != nil {
		return nil, err
	}

	//Вношу имена в мапу [user.id]username
	usernames := make(map[int]string, 4)
	for _, user := range users {
		usernames[user.Id] = user.Username
	}

	userTasksResponse := make([]domain.UserTaskResponse, 0)

	for _, task := range tasks {

		// Get team submissions for task
		submissions, err := s.taskRepo.GetTaskSubmissions(task.Id, teamId)
		if err != nil {
			return nil, err
		}

		// Build submissions response
		taskSubmissionsResponse := make([]domain.TaskSubmissionResponse, 0)
		for _, submission := range submissions {
			// Allocate sibmission object
			submissionResponse := domain.TaskSubmissionResponse{
				Username:  usernames[submission.UserId],
				Flag:      submission.Flag,
				IsCorrect: submission.IsCorrect,
				Timestamp: submission.Timestamp,
			}
			// Append allocated submission object to array of submissions
			taskSubmissionsResponse = append(taskSubmissionsResponse, submissionResponse)
		}

		// Check if task solved
		isTaskSolved, err := s.taskRepo.CheckSolvedTask(task.Id, teamId)
		if err != nil {
			return nil, err
		}

		// Allocate and fill task object
		taskResponse := domain.UserTaskResponse{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Category:    task.Description,
			Complexity:  task.Category,
			Points:      task.Points,
			Hint:        task.Hint,
			IsSolved:    isTaskSolved,
			Submissions: taskSubmissionsResponse,
			Author:      task.Author,
		}

		// Append task object to UserTasksResponse array
		userTasksResponse = append(userTasksResponse, taskResponse)
	}

	return userTasksResponse, nil
}

func (s *UserService) SolveTask(submission domain.TaskSubmission) error {
	task, err := s.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return err
	}

	if submission.Flag == task.Flag {
		submission.IsCorrect = true
	}

	if err := s.taskRepo.SolveTask(task.Id, submission.TeamId); err != nil {
		return err
	}

	if err := s.taskRepo.SaveTaskSubmission(submission); err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetTaskSubmissions(taskId, userId int) ([]*domain.TaskSubmission, error) {
	submissions, err := s.taskRepo.GetTaskSubmissions(taskId, userId)
	if err != nil {
		return nil, err
	}
	return submissions, nil
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

func (s *UserService) GetUserMe(userId int) (domain.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return domain.UserProfileResponse{}, err
	}

	team, err := s.teamRepo.GetTeamById(user.TeamId)
	if err != nil {
		return domain.UserProfileResponse{}, err
	}

	teamMembers, err := s.teamRepo.GetMembers(user.TeamId)
	if err != nil {
		return domain.UserProfileResponse{}, err
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
		Username:    user.Username,
		Contact:     user.Contact,
		TeamName:    team.Name,
		TeamMembers: membersResponse,
	}

	return userProfileResponse, nil
}

func (s *UserService) GetUser(userId int) (domain.UserProfileResponse, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return domain.UserProfileResponse{}, err
	}

	team, err := s.teamRepo.GetTeamById(user.TeamId)
	if err != nil {
		return domain.UserProfileResponse{}, err
	}

	teamMembers, err := s.teamRepo.GetMembers(user.TeamId)
	if err != nil {
		return domain.UserProfileResponse{}, err
	}

	membersResponse := make([]domain.MemberResponse, 0)

	for _, teamMember := range teamMembers {

		memberResponse := domain.MemberResponse{
			Username: teamMember.Username,
		}

		membersResponse = append(membersResponse, memberResponse)
	}

	userProfileResponse := domain.UserProfileResponse{
		Username:    user.Username,
		TeamName:    team.Name,
		TeamMembers: membersResponse,
	}

	return userProfileResponse, nil
}
