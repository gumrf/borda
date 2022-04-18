package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidate(t *testing.T) {
	type testCase struct {
		Name string

		Task SubmitFlagRequest

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.Task.Validate()

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "testValid",
		Task: SubmitFlagRequest{
			Flag: "MACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "InvalidName",
		Task: SubmitFlagRequest{
			Flag: "YACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "Invalidflag",
		Task: SubmitFlagRequest{
			Flag: "MACTF{URa_e{}o_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets1",
		Task: SubmitFlagRequest{
			Flag: "MACTF{URa_eto_zhe_fLAG213",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets2",
		Task: SubmitFlagRequest{
			Flag: "MACTFURa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})
}

func TestTaskUpdateValidate(t *testing.T) {
	type testCase struct {
		Name       string
		Prefix     string
		TaskUpdate TaskUpdate

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := tc.TaskUpdate.Validate(tc.Prefix)

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	validate(t, &testCase{
		Name: "",
		Prefix: "flag",
		TaskUpdate: TaskUpdate{
			Title: "!",
		},
		ExpectedError: nil,
	})
}

func TestUserSignInInputValidate(t *testing.T) {
	type testCase struct {
		Name string

		UserSignInInput SignInInput

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
		UserSignInInput: SignInInput{
			Username: "Jopa322",
			Password: "QAZwsx12",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TooShortUsername",
		UserSignInInput: SignInInput{
			Username: "J",
			Password: "QAZwsx12",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "ExtraSpace",
		UserSignInInput: SignInInput{
			Username: "Jopa322",
			Password: "Q fsfd",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TooShortPass",
		UserSignInInput: SignInInput{
			Username: "Jopa322",
			Password: "Q",
		},
		ExpectedError: ErrInvalidInput,
	})
}

func TestUserSignUpInputValidate(t *testing.T) {
	type testCase struct {
		Name string

		UserSignUpInput SignUpInput

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
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TestValidUUID",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "TestInvalidName",
		UserSignUpInput: SignUpInput{
			Username: "!qwewq",
			Password: "AuE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidPass",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE  322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidContact1",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "Qwe123",
			Contact:  "@drop2_der__",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidContact2",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "Qwe123",
			Contact:  "drop2_der",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "TestInvalidUUID",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidJoinTeamAttribute,
	})

	validate(t, &testCase{
		Name: "TestTooShortTeamName",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidJoinTeamAttribute,
	})

	validate(t, &testCase{
		Name: "TestInvalidTeamName",
		UserSignUpInput: SignUpInput{
			Username: "Govnoed3_3",
			Password: "AUE322%$#",
			Contact:  "@drop2_der",
		},
		ExpectedError: ErrInvalidJoinTeamAttribute,
	})
}
