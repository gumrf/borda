package domain

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidInput = errors.New("invalid input")

var ErrInvalidTeamInput = errors.New("errInvalidTeamInput")

var ErrTeamAlreadyExists = errors.New("team already exists")

var ErrTeamTokenIsInvalid = errors.New("team token is invalid or doesent exists")

var ErrTeamMembersNotFound = errors.New("team members not found")

var ErrTaskNotFound = errors.New("task not found")

var ErrTasksNotFound = errors.New("tasks not found")

var ErrTaskSolve = errors.New("can't solve task")

var ErrTaskSaveSubmission = errors.New("can't save task submission")

var ErrTaskSubmissionsNotFound = errors.New("can't get task submissions")
