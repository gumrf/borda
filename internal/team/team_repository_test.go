package team

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"borda/internal/pkg/core"
	"borda/internal/utils"
	"borda/pkg/transaction"
)

var _ core.TeamRepository = (*PostgresTeamRepository)(nil)

func TestPostgresTeamRepository_Save(t *testing.T) {
	type testCase struct {
		Name string

		Ctx  context.Context
		Team Team

		ExpectedTeam  Team
		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			db := utils.MustOpenDB(t)
			txManager := transaction.NewManager(db)
			repo := NewTeamRepository(db, txManager)

			var actualTeam Team
			var actualError error

			err := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					team, err := repo.Save(ctx, tc.Team)
					actualTeam = team

					return err
				},
				true,
			)

			actualError = err

			assert.Equal(t, tc.ExpectedTeam, actualTeam)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name: "OK",
		Ctx:  context.Background(),
		Team: Team{
			Name:  "TestTeam",
			Token: "TestToken",
		},
		ExpectedTeam: Team{
			Id:    4,
			Name:  "TestTeam",
			Token: "TestToken",
		},
		ExpectedError: nil,
	})

}

func TestPostgresTeamRepository_IsTeamFull(t *testing.T) {
	type testCase struct {
		Name string

		PostgresTeamRepository *PostgresTeamRepository

		Ctx    context.Context
		TeamId int

		ExpectedResult bool
		ExpectedError  error
	}
	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewTeamRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualResult bool
			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					actualResult = tc.PostgresTeamRepository.IsTeamFull(ctx, tc.TeamId)
					return nil
				},
				true,
			)

			assert.Equal(t, tc.ExpectedResult, actualResult)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name:                   "OK",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		TeamId:                 2,
		ExpectedResult:         false,
		ExpectedError:          nil,
	})

}

func TestPostgresTeamRepository_AddMember(t *testing.T) {
	type testCase struct {
		Name string

		PostgresTeamRepository *PostgresTeamRepository

		Ctx    context.Context
		TeamId int
		UserId int

		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewTeamRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {

			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					return tc.PostgresTeamRepository.AddMember(ctx, tc.TeamId, tc.UserId)
				},
				true,
			)

			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name:                   "Error_UserAlreadyInTeam",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		TeamId:                 2,
		UserId:                 1,
		ExpectedError:          ErrUserAlreadyInTeam,
	})

	validate(t, &testCase{
		Name:                   "Error_TeamIsFull",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		TeamId:                 1,
		UserId:                 8,
		ExpectedError:          ErrTeamIsFull,
	})

	validate(t, &testCase{
		Name:                   "OK",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		TeamId:                 3,
		UserId:                 8,
		ExpectedError:          nil,
	})

}

func TestPostgresTeamRepository_FindByToken(t *testing.T) {
	type testCase struct {
		Name string

		PostgresTeamRepository *PostgresTeamRepository

		Ctx   context.Context
		Token string

		ExpectedTeam  Team
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewTeamRepository(db, txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualTeam Team
			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					team, err := tc.PostgresTeamRepository.FindByToken(ctx, tc.Token)
					log.Println(team)
					actualTeam = team
					return err
				},
				true,
			)
			assert.Equal(t, tc.ExpectedTeam, actualTeam)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name:                   "OK",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		Token:                  "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
		ExpectedTeam: Team{
			Id:    1,
			Name:  "TestTeam1",
			Token: "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
		},
		ExpectedError: nil,
	})

	validate(t, &testCase{
		Name:                   "WrongToken",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		Token:                  "wrong-token-23434-abcd",
		ExpectedTeam:           Team{},
		ExpectedError:          ErrTeamNotFound,
	})
}

func TestPostgresTeamRepository_FindAll(t *testing.T) {
	type testCase struct {
		Name string

		PostgresTeamRepository *PostgresTeamRepository

		Ctx context.Context

		ExpectedTeams []Team
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	txManager := transaction.NewManager(db)
	repo := NewTeamRepository(utils.MustOpenDB(t), txManager)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualTeams []Team
			actualError := txManager.Run(
				tc.Ctx,
				func(ctx context.Context) error {
					teams, err := tc.PostgresTeamRepository.FindAll(ctx)
					actualTeams = teams
					return err
				},
				true,
			)
			assert.Equal(t, tc.ExpectedTeams, actualTeams)
			assert.ErrorIs(t, actualError, tc.ExpectedError)
		})
	}

	validate(t, &testCase{
		Name:                   "OK",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		ExpectedTeams: []core.Team{
			{
				Id:    1,
				Name:  "TestTeam1",
				Token: "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
			},
			{
				Id:    2,
				Name:  "TestTeam2",
				Token: "e01e949e-9dd0-428e-96c0-28adebf4df3d",
			},
			{
				Id:    3,
				Name:  "TestTeam3",
				Token: "bd58e756-7ef3-4043-bf4c-2c5ae9b9ad0b",
			},
		},
		ExpectedError: nil,
	})
}
