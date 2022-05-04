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

	type args struct {
		teamLeaderId int
		teamName     string
		username     string
		password     string
		contact      string
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
				teamName:     "Momstr",
				username:     "TestUser",
				password:     "TestPswd",
				contact:      "@testcontact",
			},
			wantResponse: 1,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			user := &domain.User{
				Username: testCase.args.username,
				Password: testCase.args.password,
				Contact:  testCase.args.contact,
			}

			helpCreateUser(t, db, user)

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

	type args struct {
		teamId   int
		username string
		password string
		contact  string
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
				teamId:   1,
				username: "User1",
				password: "User1Pswd",
				contact:  "@contact",
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
			user := &domain.User{
				Username: testCase.args.username,
				Password: testCase.args.password,
				Contact:  testCase.args.contact,
			}
			helpCreateUser(t, db, user)

			team := &domain.Team{
				Name:         testCase.wantResponse.Name,
				TeamLeaderId: testCase.wantResponse.TeamLeaderId,
			}
			helpCreateTeam(t, db, team)

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

	type args struct {
		username string
		password string
		contact  string
	}
	testTable := []struct {
		name         string
		wantResponse []*domain.Team
		wantErr      error
		args         []args
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: []args{
				{
					username: "User1",
					password: "User1Pswd",
					contact:  "@contact1",
				},
				{
					username: "User2",
					password: "User2Pswd",
					contact:  "@contact",
				},
			},
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
			for _, arg := range testCase.args {
				user := &domain.User{
					Username: arg.username,
					Password: arg.password,
					Contact:  arg.contact,
				}

				helpCreateUser(t, db, user)
			}

			for _, team := range testCase.wantResponse {
				helpCreateTeam(t, db, team)
			}

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

	type args struct {
		token        string
		username     string
		password     string
		contact      string
		teamName     string
		teamLeaderId int
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
				username:     "User1",
				password:     "User1Pswd",
				contact:      "@contact",
				teamName:     "Team1",
				teamLeaderId: 1,
			},
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
			user := &domain.User{
				Username: testCase.args.username,
				Password: testCase.args.password,
				Contact:  testCase.args.contact,
			}
			helpCreateUser(t, db, user)

			team := &domain.Team{
				Name:         testCase.args.teamName,
				TeamLeaderId: testCase.args.teamLeaderId,
			}
			helpCreateTeam(t, db, team)

			team = helpGetTeamById(t, db, testCase.wantResponse.Id)
			testCase.args.token = team.Token
			testCase.wantResponse.Token = team.Token

			actualResponse, actualErr := repo.GetTeamByToken(testCase.args.token)
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

	type args struct {
		teamId       int
		userId       int
		username1    string
		username2    string
		password1    string
		password2    string
		contact1     string
		contact2     string
		teamName     string
		teamLeaderId int
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
				teamId:       1,
				userId:       2,
				username1:    "User1",
				username2:    "User2",
				password1:    "User1Pswd",
				password2:    "User2Pswd",
				contact1:     "@contact",
				contact2:     "@contact2",
				teamName:     "Team1",
				teamLeaderId: 1,
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			user1 := &domain.User{
				Username: testCase.args.username1,
				Password: testCase.args.password1,
				Contact:  testCase.args.contact1,
			}
			helpCreateUser(t, db, user1)

			user2 := &domain.User{
				Username: testCase.args.username2,
				Password: testCase.args.password2,
				Contact:  testCase.args.contact2,
			}
			helpCreateUser(t, db, user2)

			team := &domain.Team{
				Name:         testCase.args.teamName,
				TeamLeaderId: testCase.args.teamLeaderId,
			}
			helpCreateTeam(t, db, team)

			helpSetSettings(t, db, "team_limit", "4")
			helpSetSettings(t, db, "value", "team_limit")

			actualErr := teamRepo.AddMember(testCase.args.teamId, testCase.args.userId)
			requre.Equal(testCase.wantErr, actualErr, t)
		})
	}
}

func helpCreateTeam(t *testing.T, db *sqlx.DB, team *domain.Team) *domain.Team {
	t.Helper()

	id, err := postgres.NewTeamRepository(db).SaveTeam(team.TeamLeaderId, team.Name)
	if err != nil {
		t.Fatal(err)
	}

	team.Id = id

	return team
}

func helpGetTeamById(t *testing.T, db *sqlx.DB, id int) *domain.Team {
	t.Helper()

	team, err := postgres.NewTeamRepository(db).GetTeamById(id)
	if err != nil {
		t.Fatal(err)
	}

	return team
}

func helpSetSettings(t *testing.T, db *sqlx.DB, key string, value string) int {
	t.Helper()

	id, err := postgres.NewSettingsRepository(db).Set(key, value)
	if err != nil {
		t.Fatal(err)
	}

	return id
}
