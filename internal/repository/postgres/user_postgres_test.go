package postgres_test

import (
	"borda/internal/config"
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

	hasher := hash.NewSHA1Hasher(config.PasswordSalt())

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
			name: "OK_1",
			args: args{
				username: "User2",
				password: "User2Pass",
			},
			wantResponse: &domain.User{
				Id:       2,
				Username: "User2",
				Password: "User2Pass",
				Contact:  "@user2",
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				username: "User3",
				password: "User3Pass",
			},
			wantResponse: &domain.User{
				Id:       3,
				Username: "User3",
				Password: "User3Pass",
				Contact:  "@user3",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			hashedPswd, err := hasher.Hash(testCase.args.password)
			require.Equal(testCase.wantErr, err, t)

			actualResponse, actualErr := repo.GetUserByCredentials(testCase.args.username, hashedPswd)

			hashedTestPswd, err := hasher.Hash(testCase.wantResponse.Password)
			require.Equal(testCase.wantErr, err, t)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse.Id, actualResponse.Id, t)
			require.Equal(testCase.wantResponse.Username, actualResponse.Username, t)
			require.Equal(testCase.wantResponse.Contact, actualResponse.Contact, t)
			require.Equal(hashedTestPswd, actualResponse.Password, t)
		})
	}
}

func TestUserRepository_GetUserById(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	hasher := hash.NewSHA1Hasher(config.PasswordSalt())

	type args struct {
		id int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.User
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				id: 1,
			},
			wantResponse: &domain.User{
				Id:       1,
				Username: "User1",
				Password: "User1Pass",
				Contact:  "@user1",
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				id: 3,
			},
			wantResponse: &domain.User{
				Id:       3,
				Username: "User3",
				Password: "User3Pass",
				Contact:  "@user3",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			actualResponse, actualErr := repo.GetUserById(testCase.args.id)

			hashedTestPswd, err := hasher.Hash(testCase.wantResponse.Password)
			require.Equal(testCase.wantErr, err, t)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse.Id, actualResponse.Id, t)
			require.Equal(testCase.wantResponse.Username, actualResponse.Username, t)
			require.Equal(testCase.wantResponse.Contact, actualResponse.Contact, t)
			require.Equal(hashedTestPswd, actualResponse.Password, t)
		})
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	hasher := hash.NewSHA1Hasher(config.PasswordSalt())

	testTable := []struct {
		name    string
		want    []*domain.User
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "",
			want: []*domain.User{{
				Id:       1,
				Username: "User1",
				Password: "User1Pass",
				Contact:  "@user1",
				TeamId:   1,
			},
				{
					Id:       2,
					Username: "User2",
					Password: "User2Pass",
					Contact:  "@user2",
					TeamId:   2,
				},
				{
					Id:       3,
					Username: "User3",
					Password: "User3Pass",
					Contact:  "@user3",
					TeamId:   3,
				}},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			actualResponse, actualErr := repo.GetAllUsers()
			require.Equal(testCase.wantErr, actualErr, t)

			for i, user := range testCase.want {

				wantPass, err := hasher.Hash(user.Password)
				require.Equal(testCase.wantErr, err, t)

				require.Equal(user.Id, actualResponse[i].Id, t)
				require.Equal(wantPass, actualResponse[i].Password, t)
				require.Equal(user.Contact, actualResponse[i].Contact, t)
				require.Equal(user.Username, actualResponse[i].Username, t)
				require.Equal(user.TeamId, actualResponse[i].TeamId)

			}

		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	hasher := hash.NewSHA1Hasher(config.PasswordSalt())

	type args struct {
		userId      int
		newPassword string
	}
	testTable := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				userId:      1,
				newPassword: "NewPassword",
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				userId:      2,
				newPassword: "NewPswd",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			pswd, err := hasher.Hash(testCase.args.newPassword)
			require.Equal(testCase.wantErr, err, t)

			actualErr := repo.UpdatePassword(testCase.args.userId, pswd)
			require.Equal(testCase.wantErr, actualErr, t)
		})
	}
}

//Не используется функция
// func Test_AssignRole(t *testing.T) {
// 	db := MustOpenDB(t)
// 	repo := postgres.NewUserRepository(db)
// 	require := require.New(t)

// 	type args struct {
// 		userId int
// 		roleId int
// 	}
// 	testTable := []struct {
// 		name    string
// 		args    args
// 		wantErr error
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "OK",
// 			args: args{
// 				userId: 4,
// 				roleId: 2,
// 			},
// 			wantErr: nil,
// 		},
// 	}
// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			TestUserRepository_SaveUser(t)

// 			actualErr := repo.AssignRole(testCase.args.userId, testCase.args.roleId)
// 			require.Equal(testCase.wantErr, actualErr, t)
// 		})
// 	}
// }

func TestUserRepository_GetUserRole(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewUserRepository(db)
	require := require.New(t)

	type args struct {
		userId int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.Role
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				userId: 1,
			},
			wantResponse: &domain.Role{
				Id:   1,
				Name: "admin",
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				userId: 2,
			},
			wantResponse: &domain.Role{
				Id:   2,
				Name: "user",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := repo.GetUserRole(testCase.args.userId)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}
