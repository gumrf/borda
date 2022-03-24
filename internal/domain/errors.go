package domain

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")

var ErrInvalidInput = errors.New("invalid input")

var ErrInvalidTeamInput = errors.New("ErrInvalidTeamInput")
