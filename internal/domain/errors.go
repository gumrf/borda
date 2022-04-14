package domain

import "errors"

type ErrorResponse struct {
	Status string `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
}

var (
	ErrInvalidInput = errors.New("invalid input")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUsersNotFound     = errors.New("users not found")

	ErrTeamAlreadyExists   = errors.New("team already exists")
	ErrInvalidTeamInput    = errors.New("errInvalidTeamInput")
	ErrTeamTokenIsInvalid  = errors.New("team token is invalid or doesent exists")
	ErrTeamNotFound        = errors.New("team not found")
	ErrTeamMembersNotFound = errors.New("team members not found")
	ErrMembers             = errors.New("members not found")

	ErrTaskCreate              = errors.New("can't create new task")
	ErrTaskUpdate              = errors.New("can't update task")
	ErrTasksNotFound           = errors.New("tasks not found")
	ErrTaskNotFound            = errors.New("task not found")
	ErrTaskSolve               = errors.New("can't solve task")
	ErrTaskSaveSubmission      = errors.New("can't save task submission")
	ErrTaskSubmissionsNotFound = errors.New("can't get task submissions")
)
