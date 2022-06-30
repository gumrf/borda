package team

import (
	"borda/internal/pkg/core"
	"borda/internal/utils"
	"borda/pkg/transaction"

	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ core.TeamRepository = (*PostgresTeamRepository)(nil)

func TestPostgresTeamRepository_Save(t *testing.T) {
	type testCase struct {
		Name string

		PostgresTeamRepository *PostgresTeamRepository

		Ctx  context.Context
		Team Team

		ExpectedTeam  Team
		ExpectedError error
	}
	db := utils.MustOpenDB(t)
	repo := NewTeamRepository(db)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualTeam Team
			var actualError error

			err := transaction.Transactional(
				tc.Ctx, db,
				func(ctx context.Context) error {
					_team, err := tc.PostgresTeamRepository.Save(ctx, tc.Team)
					actualTeam = _team

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
		Name:                   "OK",
		PostgresTeamRepository: repo,
		Ctx:                    context.Background(),
		Team: Team{
			Name:         "TestTeam",
			TeamLeaderId: 1,
			Token:        "TestToken",
		},
		ExpectedTeam: Team{
			Id:           4,
			Name:         "TestTeam",
			TeamLeaderId: 1,
			Token:        "TestToken",
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
	repo := NewTeamRepository(db)

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var actualResult bool
			actualError := transaction.Transactional(tc.Ctx, db,
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

		// ExpectedTeam  Team
		ExpectedError error
	}

	db := utils.MustOpenDB(t)
	repo := NewTeamRepository(utils.MustOpenDB(t))

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {

			actualError := transaction.Transactional(tc.Ctx, db,
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
