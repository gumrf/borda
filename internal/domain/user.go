package domain

import (
	"regexp"
	
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"name"`
	Password string `json:"password" db:"password"`
	Contact  string `json:"contact" db:"contact"`
	TeamId   int    `json:"teamId" db:"team_id"`
}

type UserSignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`

	// ?team[create, join]=[teamName, token]

	AttachTeamMethod    string `json:"attachTeamMethod"`
	AttachTeamAttribute string `json:"attachTeamAttribute"`
}

type UserSignInInput struct {
	Username string
	Password string
}

func (t UserSignInInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 128), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),
	)

	if err != nil{
		return ErrInvalidInput
	}
	
	return nil
}

func (t UserSignUpInput) Validate() error {
	err := validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 50), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(8, 100), validation.Match(regexp.MustCompile("^[0-9a-zA-Z!@#$%^&*]+$"))),

		validation.Field(&t.Contact, validation.Required, validation.Length(5, 50), validation.Match(regexp.MustCompile("^@[0-9_a-z]+[^_]$"))),
	)
	
	if err != nil{
		return ErrInvalidInput
	}
	
	return nil

}
