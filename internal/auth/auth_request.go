package auth

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
}

func (rr RegistrationRequest) Validate() error {
	err := validation.ValidateStruct(&rr,
		validation.Field(&rr.Username,
			validation.Required,
			validation.Length(2, 50),
			validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$")),
		),
		validation.Field(&rr.Password,
			validation.Required,
			validation.Length(8, 100),
			validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$")),
		),
		validation.Field(&rr.Contact,
			validation.Required,
			validation.Length(5, 100),
			validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$")),
		),
	)

	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

type AuthorizationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ar AuthorizationRequest) Validate() error {
	err := validation.ValidateStruct(&ar,
		validation.Field(&ar.Username,
			validation.Required,
			validation.Length(2, 50),
			validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$")),
		),
		validation.Field(&ar.Password,
			validation.Required,
			validation.Length(8, 128),
			validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$")),
		),
	)

	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

type ConfirmEmailRequest struct {
	Code string `json:"code"`
}
