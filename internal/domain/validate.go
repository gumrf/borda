package domain

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (t UserSignInInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 128), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (t UserSignUpInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 100), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),

		validation.Field(&t.Contact, validation.Required, validation.Length(5, 100), validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	//if t.AttachTeamMethod == "create" {
	//	err := validation.Validate(&t.AttachTeamAttribute, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$")))
	//	if err != nil {
	//		return ErrInvalidTeamInput
	//	}
	//} else if t.AttachTeamMethod == "join" {
	//	err := validation.Validate(&t.AttachTeamAttribute, validation.Required, is.UUIDv4)
	//	if err != nil {
	//		return ErrInvalidTeamInput
	//	}
	//} else {
	//	return ErrInvalidTeamInput
	//}

	return nil

}

func (t Task) Validate() error {
	err := validation.ValidateStruct(&t,
		validation.Field(&t.Id, validation.Required, is.Digit),
		validation.Field(&t.Title, validation.Required, validation.Match(regexp.MustCompile("^[0-9A-Za-z_?!,.\\s]+$"))),
		validation.Field(&t.Description, validation.Required),
		validation.Field(&t.Category, validation.Required, is.LowerCase),
		validation.Field(&t.Complexity, validation.Required, is.LowerCase),
		validation.Field(&t.Points, validation.Required, is.Digit),
		validation.Field(&t.Hint),
		validation.Field(&t.Flag, validation.Required, validation.Match(regexp.MustCompile("^MACTF{[0-9A-Za-z_]+}$"))),
		validation.Field(&t.IsActive, validation.Required),
		validation.Field(&t.IsDisabled, validation.Required),
		validation.Field(&t.Author),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (t TaskUpdate) Validate() error {
	err := validation.ValidateStruct(&t,
		validation.Field(&t.Title, validation.Match(regexp.MustCompile("^[0-9A-Za-z_?!,.\\s]+$"))),
		validation.Field(&t.Description),
		validation.Field(&t.Category, is.LowerCase),
		validation.Field(&t.Complexity, is.LowerCase),
		validation.Field(&t.Points, is.Digit),
		validation.Field(&t.Hint),
		validation.Field(&t.Flag, validation.Match(regexp.MustCompile("^MACTF{[0-9A-Za-z_]+}$"))),
		validation.Field(&t.AuthorName, validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		validation.Field(&t.AuthorContact, validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (t SubmitTaskRequest) Validate() error {
	err := validation.ValidateStruct(&t,
		validation.Field(&t.Flag, validation.Required, validation.Match(regexp.MustCompile("^MACTF{[0-9A-Za-z_]+}$"))),
	)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func (a Author) Validate() error {
	err := validation.ValidateStruct(&a,
		validation.Field(&a.Id, is.Digit),
		validation.Field(&a.Name, validation.Required, validation.Match(regexp.MustCompile("^[0-9A-Za-z_?!,.\\s]+$"))),
		validation.Field(&a.Contact, validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
