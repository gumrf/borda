package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	hash "borda/pkg/hash"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserRepository_SaveUser(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	type args struct {
		username string
		password string
		contact  string
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse int
		wantErr      error
	}{
		// TODO: Add test casese
		{
			name: "OK",
			args: args{
				username: "Success",
				password: "Sucsess",
				contact:  "@success",
			},
			wantResponse: 4,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := repo.SaveUser(testCase.args.username, testCase.args.password, testCase.args.contact)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestUserRepository_GetUserByCredentials(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	var hasher hash.PasswordHasher

	type args struct {
		username string
		password string
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.User
		wantErr      error
		hasher       hash.PasswordHasher
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				username: "User2",
				password: "User2Pass",
			},
			wantResponse: &domain.User{
				Id:       2,
				Username: "User2",
				Password: "User2Pass",
				Contact:  "@user2",
				TeamId:   2,
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			hashedPswd, err := hasher.Hash(testCase.args.password)
			require.Equal(testCase.wantErr, err, t)

			actualResponse, actualErr := repo.GetUserByCredentials(testCase.args.username, hashedPswd)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}
