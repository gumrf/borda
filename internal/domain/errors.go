package domain

import "errors"

var ErrUserAlreadyExists = errors.New("User already exists")

var ErrInvalidInput = errors.New("Invalid input")

var ErrInvalidTeamInput = errors.New("ErrInvalidTeamInput")

var ErrTeamAlreadyExists = errors.New("Team already exists")

var ErrTeamTokenIsInvalid = errors.New("Team token is invalid or doesent exists")
