package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSignInInputValidate(t *testing.T) {
	type testCase struct {
		Name string

		UserSignInInput UserSignInInput

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.UserSignInInput.Validate()

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "TestValid",
		UserSignInInput: UserSignInInput{
			Username: "Jopa322",
			Password: "QAZwsx12",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TooShortUsername",
		UserSignInInput: UserSignInInput{
			Username: "J",
			Password: "QAZwsx12",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "ExtraSpace",
		UserSignInInput: UserSignInInput{
			Username: "Jopa322",
			Password: "Q fsfd",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TooShortPass",
		UserSignInInput: UserSignInInput{
			Username: "Jopa322",
			Password: "Q",
		},
		ExpectedError: ErrInvalidInput,
	})
}

func TestUserSignUpInputValidate(t *testing.T) {
	type testCase struct {
		Name string

		UserSignUpInput UserSignUpInput

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.UserSignUpInput.Validate()

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "TestValidCreate",
		UserSignUpInput: UserSignUpInput{
			Username:            "Govnoed3_3",
			Password:            "AUE322%$#",
			Contact:             "@drop2_der",
			AttachTeamMethod:    "create",
			AttachTeamAttribute: "DT_404",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TestValidUUID",
		UserSignUpInput: UserSignUpInput{
			Username:            "Govnoed3_3",
			Password:            "AUE322%$#",
			Contact:             "@drop2_der",
			AttachTeamMethod:    "join",
			AttachTeamAttribute: "6e3ab84d-1ca8-4bf5-a803-123c66f965ae",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TestInvalidName",
		UserSignUpInput: UserSignUpInput{
			Username: "!qwewq",
			Password: "AuE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidPass",
		UserSignUpInput: UserSignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE  322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidContact1",
		UserSignUpInput: UserSignUpInput{
			Username: "Govnoed3_3",
			Password: "Qwe123",
			Contact:  "@drop2_der__",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidContact2",
		UserSignUpInput: UserSignUpInput{
			Username: "Govnoed3_3",
			Password: "Qwe123",
			Contact:  "drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidUUID",
		UserSignUpInput: UserSignUpInput{
			Username:            "Govnoed3_3",
			Password:            "AUE322%$#",
			Contact:             "@drop2_der",
			AttachTeamMethod:    "join",
			AttachTeamAttribute: "6e3ab84f965ae",
		},
		ExpectedError: ErrInvalidTeamInput,
	})

	validate(t, &testCase{
		Name: "TestTooShortTeamName",
		UserSignUpInput: UserSignUpInput{
			Username:            "Govnoed3_3",
			Password:            "AUE322%$#",
			Contact:             "@drop2_der",
			AttachTeamMethod:    "create",
			AttachTeamAttribute: "a",
		},
		ExpectedError: ErrInvalidTeamInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidTeamName",
		UserSignUpInput: UserSignUpInput{
			Username:            "Govnoed3_3",
			Password:            "AUE322%$#",
			Contact:             "@drop2_der",
			AttachTeamMethod:    "create",
			AttachTeamAttribute: "!DT_404",
		},
		ExpectedError: ErrInvalidTeamInput,
	})
}
