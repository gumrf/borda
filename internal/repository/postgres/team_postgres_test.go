package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestTeamRepository_SaveTeam(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	helpCreateUser(t, db)

	type args struct {
		teamLeaderId int
		teamName     string
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
				teamLeaderId: 1,
				teamName:     "Team1",
			},
			wantResponse: 1,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := teamRepo.SaveTeam(testCase.args.teamLeaderId, testCase.args.teamName)

			requre.Equal(testCase.wantErr, actualErr, t)
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTeamRepository_GetTeamById(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	helpCreateUser(t, db)
	helpCreateTeam(t, db)

	type args struct {
		teamId int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.Team
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				teamId: 1,
			},
			wantResponse: &domain.Team{
				Id:           1,
				Name:         "Team1",
				TeamLeaderId: 1,
				Members: []domain.TeamMember{
					{
						UserId: 1,
						Name:   "User1",
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := teamRepo.GetTeamById(testCase.args.teamId)

			requre.Equal(testCase.wantErr, actualErr, t)
			requre.Equal(testCase.wantResponse.Id, actualResponse.Id, t)
			requre.Equal(testCase.wantResponse.Members, actualResponse.Members, t)
			requre.Equal(testCase.wantResponse.Name, actualResponse.Name, t)
			requre.Equal(testCase.wantResponse.TeamLeaderId, actualResponse.TeamLeaderId, t)
		})
	}
}

func TestTeamRepository_GetTeams(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	helpCreateUser(t, db)
	helpCreateTeam(t, db)

	testTable := []struct {
		name         string
		wantResponse []*domain.Team
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			wantResponse: []*domain.Team{
				{
					Id:           1,
					Name:         "Team1",
					TeamLeaderId: 1,
					Members: []domain.TeamMember{
						{
							UserId: 1,
							Name:   "User1",
						},
					},
				},
				{
					Id:           2,
					Name:         "Team2",
					TeamLeaderId: 2,
					Members: []domain.TeamMember{
						{
							UserId: 2,
							Name:   "User2",
						},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := teamRepo.GetTeams()
			requre.Equal(testCase.wantErr, actualErr, t)

			for i, testTeam := range testCase.wantResponse {
				requre.Equal(testTeam.Id, actualResponse[i].Id, t)
				requre.Equal(testTeam.Name, actualResponse[i].Name, t)
				requre.Equal(testTeam.TeamLeaderId, actualResponse[i].TeamLeaderId, t)
				requre.Equal(testTeam.Members, actualResponse[i].Members, t)
			}
		})
	}
}

func TestTeamRepository_GetTeamByToken(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	helpCreateUser(t, db)
	helpCreateTeam(t, db)

	testTable := []struct {
		name         string
		wantResponse *domain.Team
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			wantResponse: &domain.Team{
				Id:           1,
				Name:         "Team1",
				TeamLeaderId: 1,
				Token:        "",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			team := helpGetTeamById(t, db, testCase.wantResponse.Id)
			testCase.wantResponse.Token = team.Token

			actualResponse, actualErr := repo.GetTeamByToken(testCase.wantResponse.Token)
			requre.Equal(testCase.wantErr, actualErr, t)
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

//Add: team tessting error response!
func TestTeamRepository_AddMember(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	helpCreateUser(t, db)
	helpCreateTeam(t, db)

	type args struct {
		teamId int
		userId int
	}
	testTable := []struct {
		name string

		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				teamId: 1,
				userId: 3,
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			helpSetSettings(t, db, "team_limit", "4")
			helpSetSettings(t, db, "value", "team_limit")

			actualErr := teamRepo.AddMember(testCase.args.teamId, testCase.args.userId)
			requre.Equal(testCase.wantErr, actualErr, t)
		})
	}
}

// Вызывать только после helpCreateUser()!!!!!
func helpCreateTeam(t *testing.T, db *sqlx.DB) {
	t.Helper()

	teams := []*domain.Team{
		{
			Name:         "Team1",
			TeamLeaderId: 1,
		},
		{
			Name:         "Team2",
			TeamLeaderId: 2,
		},
	}

	for _, team := range teams {
		id, err := postgres.NewTeamRepository(db).SaveTeam(team.TeamLeaderId, team.Name)
		if err != nil {
			t.Fatal(err)
		}

		team.Id = id
	}
}

func helpGetTeamById(t *testing.T, db *sqlx.DB, id int) *domain.Team {
	t.Helper()

	team, err := postgres.NewTeamRepository(db).GetTeamById(id)
	if err != nil {
		t.Fatal(err)
	}

	return team
}
