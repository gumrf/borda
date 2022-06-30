package domain

// import (
// 	"fmt"
// 	"regexp"

// 	validation "github.com/go-ozzo/ozzo-validation/v4"
// 	"github.com/go-ozzo/ozzo-validation/v4/is"
// )


// func (t Task) Validate(prefix string) error {
// 	err := validation.ValidateStruct(&t,
// 		validation.Field(&t.Title, validation.Required, validation.Match(regexp.MustCompile("^[A-Za-z0-9 ]+$"))),
// 		validation.Field(&t.Description, validation.Required),
// 		validation.Field(&t.Category, validation.Required, is.LowerCase),
// 		validation.Field(&t.Complexity, validation.Required, is.LowerCase),
// 		validation.Field(&t.Hint),
// 		validation.Field(&t.Flag, validation.Match(regexp.MustCompile(fmt.Sprintf("^%s{[0-9A-Za-z_]+}$", prefix)))),
// 		validation.Field(&t.IsActive, validation.Required),
// 		validation.Field(&t.IsDisabled, validation.Required),
// 		validation.Field(&t.Author),
// 	)

// 	if err != nil {
// 		return fmt.Errorf("validation error: %v", err)
// 	}

// 	if t.Points != 0 {
// 		if t.Points < 0{
// 			return fmt.Errorf("validation error: points must be > 0")
// 		}
// 	}

// 	return nil
// }

// func (u TaskUpdate) Validate(prefix string) error {
// 	err := validation.ValidateStruct(&u,
// 		validation.Field(&u.Title, validation.Match(regexp.MustCompile("^[A-Za-z0-9 ]+$"))),
// 		validation.Field(&u.Description),
// 		validation.Field(&u.Category, is.LowerCase),
// 		validation.Field(&u.Complexity, is.LowerCase),
// 		validation.Field(&u.Hint),
// 		validation.Field(&u.Flag, validation.Match(regexp.MustCompile(fmt.Sprintf("^%s{[0-9A-Za-z_]+}$", prefix)))),
// 		validation.Field(&u.AuthorName, validation.Match(regexp.MustCompile("^[0-9A-Za-z_]+$"))),
// 		validation.Field(&u.AuthorContact, validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
// 	)

// 	if err != nil {
// 		return fmt.Errorf("validation error: %v", err)
// 	}

// 	if u.Points != 0 {
// 		if u.Points < 0{
// 			return fmt.Errorf("validation error: points must be > 0")
// 		}
// 	}

// 	return nil
// }

// func (t SubmitFlagRequest) Validate() error {
// 	err := validation.ValidateStruct(&t,
// 		validation.Field(&t.Flag, validation.Required, validation.Match(regexp.MustCompile("^flag{[0-9A-Za-z_]+}$"))),
// 	)
// 	if err != nil {
// 		return fmt.Errorf("validation error: %v", err)
// 	}

// 	return nil
// }

// func (a Author) Validate() error {
// 	err := validation.ValidateStruct(&a,
// 		validation.Field(&a.Name, validation.Required, validation.Match(regexp.MustCompile("^[0-9A-Za-z_?!,.\\s]+$"))),
// 		validation.Field(&a.Contact, validation.Match(regexp.MustCompile("^@[a-z0-9_]+[a-z0-9]$"))),
// 	)

// 	if err != nil {
// 		return fmt.Errorf("validation error: %v", err)
// 	}

// 	// if a.Id <= 0 {
// 	// 	return fmt.Errorf("validation err: id must be > 0")
// 	// }

// 	return nil
// }
