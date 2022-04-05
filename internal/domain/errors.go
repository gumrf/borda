package domain

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidInput = errors.New("invalid input")

var ErrInvalidTeamInput = errors.New("errInvalidTeamInput")

var ErrTeamAlreadyExists = errors.New("team already exists")

var ErrTeamTokenIsInvalid = errors.New("team token is invalid or doesent exists")

var ErrUserNotFound = errors.New("user not found")

var ErrUsersNotFound = errors.New("users not found")

var ErrTeamNotFound = errors.New("team not found")

var ErrMembers = errors.New("members not found")
