package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidate(t *testing.T) {
	type testCase struct {
		Name string

		Task SubmitTaskRequest

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
		Task: SubmitTaskRequest{
			Flag: "MACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "InvalidName",
		Task: SubmitTaskRequest{
			Flag: "YACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "Invalidflag",
		Task: SubmitTaskRequest{
			Flag: "MACTF{URa_e{}o_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets1",
		Task: SubmitTaskRequest{
			Flag: "MACTF{URa_eto_zhe_fLAG213",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets2",
		Task: SubmitTaskRequest{
			Flag: "MACTFURa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})
}
