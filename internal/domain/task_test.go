package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskValidate(t *testing.T) {
	type testCase struct {
		Name string

		Task Task

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
		Task: Task{
			Flag: "MACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "InvalidName",
		Task: Task{
			Flag: "YACTF{URa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "Invalidflag",
		Task: Task{
			Flag: "MACTF{URa_e{}o_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets1",
		Task: Task{
			Flag: "MACTF{URa_eto_zhe_fLAG213",
		},
		ExpectedError: ErrInvalidInput,
	})

	validate(t, &testCase{
		Name: "InvalidBrackets2",
		Task: Task{
			Flag: "MACTFURa_eto_zhe_fLAG213}",
		},
		ExpectedError: ErrInvalidInput,
	})
}
