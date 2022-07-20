package user

import (
	"borda/internal/pkg/core"
	"borda/internal/utils"
	"borda/pkg/transaction"

	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ core.UserRepository = (*PostgresUserRepository)(nil)

func TestPostgresUserRepository_Save(t *testing.T) {
	type testCase struct {
		Name string

		Ctx  context.Context
		User User

		ExpectedUser  User
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewUserRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			//actualUser, actualError := tc.PostgresUserRepository.Save(tc.Ctx, tc.User)

			var actualUser User

			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					user, err := repo.Save(ctx, tc.User)
					actualUser = user
					return err
				},
				true,
			)

			assert.Equal(t, tc.ExpectedUser, actualUser)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Ctx:  context.Background(),
		User: User{
			Id:       0,
			Username: "Max",
			Password: "m4x-1s-g0d",
			Contact:  "max@mail.god",
			TeamId:   0,
		},
		ExpectedUser: User{
			Id:       1,
			Username: "Max",
			Password: "m4x-1s-g0d",
			Contact:  "max@mail.god",
			TeamId:   0,
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name: "DuplicateName",
		Ctx:  context.Background(),
		User: User{
			Id:       0,
			Username: "Max",
			Password: "m4x-1s-g0d",
			Contact:  "max@mail.god",
			TeamId:   0,
		},
		ExpectedUser:  User{},
		ExpectedError: ErrUserAlreadyExists,
	})
}

func TestPostgresUserRepository_SaveAll(t *testing.T) {
	type testCase struct {
		Name string

		Ctx      context.Context
		Entities []User

		ExpectedSlice []User
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewUserRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualUsers []User
			err := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					users, err := repo.SaveAll(ctx, tc.Entities)
					actualUsers = users

					return err
				},
				true,
			)

			actualError := err

			assert.Equal(t, tc.ExpectedSlice, actualUsers)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Ctx:  context.Background(),
		Entities: []User{
			{
				Id:       1,
				Username: "Max",
				Password: "m4x-1s-g0d",
				Contact:  "max@mail.god",
				TeamId:   0,
			},
			{
				Id:       2,
				Username: "Max2",
				Password: "m4x-1s-g0d-2",
				Contact:  "max2@mail.god",
				TeamId:   0,
			},
		},
		ExpectedSlice: []User{
			{
				Id:       1,
				Username: "Max",
				Password: "m4x-1s-g0d",
				Contact:  "max@mail.god",
				TeamId:   0,
			},
			{
				Id:       2,
				Username: "Max2",
				Password: "m4x-1s-g0d-2",
				Contact:  "max2@mail.god",
				TeamId:   0,
			},
		},
		ExpectedError: nil,
	})
}

func TestPostgresUserRepository_FindById(t *testing.T) {
	type testCase struct {
		Name string

		Ctx context.Context
		Id  int

		ExpectedUser  User
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewUserRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			_, _ = repo.Save(
				tc.Ctx,
				User{
					Username: "Max",
					Password: "m4x-1s-g0d",
					Contact:  "max@mail.god",
				},
			)

			var actualUser User

			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					user, err := repo.FindById(ctx, tc.Id)
					actualUser = user
					return err
				},
				true,
			)

			assert.Equal(t, tc.ExpectedUser, actualUser)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Ctx:  context.Background(),
		Id:   1,
		ExpectedUser: User{
			Id:       1,
			Username: "Max",
			Password: "m4x-1s-g0d",
			Contact:  "max@mail.god",
			TeamId:   0,
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name:          "UserNotFound",
		Ctx:           context.Background(),
		Id:            2,
		ExpectedUser:  User{},
		ExpectedError: ErrUserNotFound,
	})

	validate(t, &testCase{
		Name:          "NegativeId",
		Ctx:           context.Background(),
		Id:            -2,
		ExpectedUser:  User{},
		ExpectedError: ErrUserNotFound,
	})
}

// func TestPostgresUserRepository_FindByCredentials(t *testing.T) {
// 	type testCase struct {
// 		Name string

// 		PostgresUserRepository *PostgresUserRepository

// 		Username string
// 		Password string

// 		ExpectedUser  User
// 		ExpectedError error
// 	}

// 	validate := func(t *testing.T, tc *testCase) {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			_, _ = tc.PostgresUserRepository.Save(
// 				context.Background(),
// 				User{
// 					Username: "Max",
// 					Password: "m4x-1s-g0d",
// 					Contact:  "max@mail.god",
// 				},
// 			)

// 			actualUser, actualError := tc.PostgresUserRepository.FindByCredentials(tc.Username, tc.Password)

// 			assert.Equal(t, tc.ExpectedUser, actualUser)
// 			assert.ErrorIs(t, actualError, tc.ExpectedError)
// 		})
// 	}

// 	validate(t, &testCase{
// 		Name:                   "OK",
// 		PostgresUserRepository: NewUserRepository(utils.MustOpenDB(t)),
// 		Username:               "Max",
// 		Password:               "m4x-1s-g0d",
// 		ExpectedUser: User{
// 			Id:       1,
// 			Username: "Max",
// 			Password: "m4x-1s-g0d",
// 			Contact:  "max@mail.god",
// 			TeamId:   0,
// 		},
// 		ExpectedError: nil,
// 	})

// 	validate(t, &testCase{
// 		Name:                   "NotFound",
// 		PostgresUserRepository: NewUserRepository(utils.MustOpenDB(t)),
// 		Username:               "NeMax",
// 		Password:               "m4x-1s-g0d",
// 		ExpectedUser:           User{},
// 		ExpectedError:          ErrUserNotFound,
// 	})
// }
