package postgres_test

import (
	"borda/internal/repository/postgres"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestSettingsRepository_Get(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewSettingsRepository(db)
	require := require.New(t)

	type args struct {
		key string
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse string
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				key: "team_limit",
			},
			wantResponse: "4",
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			helpSetSettings(t, db, "team_limit", "4")

			actualResponse, actualErr := repo.Get(testCase.args.key)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestSettingsRepository_Set(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewSettingsRepository(db)
	require := require.New(t)

	type args struct {
		key   string
		value string
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse int
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				key:   "TEST",
				value: "1337",
			},
			wantResponse: 1,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := repo.Set(testCase.args.key, testCase.args.value)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func helpSetSettings(t *testing.T, db *sqlx.DB, key string, value string) int {
	t.Helper()

	id, err := postgres.NewSettingsRepository(db).Set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	return id
}
