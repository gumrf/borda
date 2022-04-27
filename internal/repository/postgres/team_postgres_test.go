package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTeamRepository_SaveTeam(t *testing.T) {
	db := MustOpenDB(t)
	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

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
				teamLeaderId: 4,
				teamName:     "Momstr",
			},
			wantResponse: 4,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := userRepo.SaveUser("TestUser", "TestPswd", "@testcontact")
			requre.Equal(testCase.wantErr, err, t)

			actualResponse, actualErr := teamRepo.SaveTeam(testCase.args.teamLeaderId, testCase.args.teamName)

			requre.Equal(testCase.wantErr, actualErr, t)
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTeamRepository_GetTeamById(t *testing.T) {
	db := MustOpenDB(t)
	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

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
				Token:        "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
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
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTeamRepository_GetTeams(t *testing.T) {
	db := MustOpenDB(t)
	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

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
					Token:        "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
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
					Token:        "e01e949e-9dd0-428e-96c0-28adebf4df3d",
					Members: []domain.TeamMember{
						{
							UserId: 2,
							Name:   "User2",
						},
					},
				},
				{
					Id:           3,
					Name:         "Team3",
					TeamLeaderId: 3,
					Token:        "bd58e756-7ef3-4043-bf4c-2c5ae9b9ad0b",
					Members: []domain.TeamMember{
						{
							UserId: 3,
							Name:   "User3",
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
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTeamRepository_GetTeamByToken(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewTeamRepository(db)
	requre := require.New(t)

	type args struct {
		token string
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.Team
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				token: "e01e949e-9dd0-428e-96c0-28adebf4df3d",
			},
			wantResponse: &domain.Team{
				Id:           2,
				Name:         "Team2",
				TeamLeaderId: 2,
				Token:        "e01e949e-9dd0-428e-96c0-28adebf4df3d",
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				token: "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
			},
			wantResponse: &domain.Team{
				Id:           1,
				Name:         "Team1",
				TeamLeaderId: 1,
				Token:        "7c8b4c73-9fdd-4fbb-b926-43de9aa6f24d",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := repo.GetTeamByToken(testCase.args.token)

			requre.Equal(testCase.wantErr, actualErr, t)
			requre.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTeamRepository_AddMember(t *testing.T) {
	db := MustOpenDB(t)
	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	requre := require.New(t)

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
				userId: 4,
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := userRepo.SaveUser("TestUser", "TestPswd", "@testcontact")
			requre.Equal(testCase.wantErr, err, t)

			actualErr := teamRepo.AddMember(testCase.args.teamId, testCase.args.userId)
			requre.Equal(testCase.wantErr, actualErr, t)
		})
	}
}
