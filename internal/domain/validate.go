package domain

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (t SignInInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 128), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

func (t SignUpInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 100), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),

		validation.Field(&t.Contact, validation.Required, validation.Length(5, 100), validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

func (t TeamInput) Validate() error {
	if t.Method == "create" {
		if err := validation.Validate(&t.Attribute, validation.Required, validation.Length(3, 50),
			validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))); err != nil {
			return ErrInvalidTeamInput
		}
	} else if t.Method == "join" {
		if err := validation.Validate(&t.Attribute, validation.Required, is.UUIDv4); err != nil {
			return ErrInvalidTeamInput
		}
	} else {
		return ErrInvalidMethod
	}

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
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

func (u TaskUpdate) Validate() error {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Title, validation.Match(regexp.MustCompile("^[0-9A-Za-z_?!,.\\s]+$"))),
		validation.Field(&u.Description),
		validation.Field(&u.Category, is.LowerCase),
		validation.Field(&u.Complexity, is.LowerCase),
		//validation.Field(&u.Points, is.Digit), validation error: points: must be either a string or byte slice.
		validation.Field(&u.Hint),
		validation.Field(&u.Flag, validation.Match(regexp.MustCompile("^MACTF{[0-9A-Za-z_]+}$"))),
		validation.Field(&u.AuthorName, validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		validation.Field(&u.AuthorContact, validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
	)

	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

func (t SubmitFlagRequest) Validate() error {
	err := validation.ValidateStruct(&t,
		validation.Field(&t.Flag, validation.Required, validation.Match(regexp.MustCompile("^flag{[0-9A-Za-z_]+}$"))),
	)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
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
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}
