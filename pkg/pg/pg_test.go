package pg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {

	type testCase struct {
		Name string

		Url string

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			_, actualError := Connect(tc.Url)

			require.NoError(t, actualError)
			// assert.ErrorIs(t, actualError, tc.ExpectedError)

		})
	}

	validate(t, &testCase{
		Name:          "OK",
		Url:           "postgres://postgres:secret_password@127.0.0.1:5432/borda?sslmode=disable",
		ExpectedError: nil,
	})

}

func TestMigrate(t *testing.T) {

	type testCase struct {
		Name string

		ConnectionURL string
		Source        string

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := Migrate(tc.ConnectionURL, tc.Source, 2)

			require.NoError(t, actualError)
			// assert.ErrorIs(t, actualError, tc.ExpectedError)

		})
	}

	validate(t, &testCase{
		Name:          "OK",
		ConnectionURL: "pgx://postgres:secret_password@127.0.0.1:5432/borda?sslmode=disable",
		Source:        "file://../../migrations",
		ExpectedError: nil,
	})

}
