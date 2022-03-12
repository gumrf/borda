package domain

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"name"`
	Password string `json:"password" db:"password"`
	Contact  string `json:"contact" db:"contact"`
	TeamId   int    `json:"teamId" db:"team_id"`
}

<<<<<<< HEAD
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

func (t *User) validate() error {
=======
func (t *User) Validate() error {
>>>>>>> b9359c4 (fix user validate)
	return validation.ValidateStruct(&t,
		// Username cannot be empty, and the length must between 2 and 20, may contains letters, numbers and '_'
		validation.Field(&t.Username, validation.Required, validation.Length(2, 20), validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
		// Password cannot be empty, and the length must between 4 and 100, and must contain Uppercase letter, lowcase letter, and numbers
		validation.Field(&t.Password, validation.Required, validation.Length(4, 100), is.Digit, is.LowerCase, is.UpperCase, validation.Match(regexp.MustCompile("^[^ ]+$"))),
	)

}
