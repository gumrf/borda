package team

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var (
	ErrInvalidJoinTeamMethod    = errors.New("invalid join team method")
	ErrInvalidJoinTeamAttribute = errors.New("invalid join team attribute")
)

type CreatTeamRequest struct {
	Name string `json:"name"`
}

type JoinTeamRequest struct {
	Token string `json:"token"`
}

func (ctr CreatTeamRequest) Validate() error {
	err := validation.Validate(
		&ctr.Name,
		validation.Required,
		validation.Length(3, 50),
		validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$")),
	)
	if err != nil {
		return ErrInvalidJoinTeamAttribute
	}

	return nil
}

func (jtr JoinTeamRequest) Validate() error {
	err := validation.Validate(
		&jtr.Token,
		validation.Required,
		is.UUIDv4,
	)
	if err != nil {
		return ErrInvalidJoinTeamAttribute
	}

	return nil
}
