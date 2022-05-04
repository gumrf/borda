package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_SaveTask(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		input domain.Task
	}

	testTable := []struct {
		name         string
		args         args
		wantResponse int
		wantErr      error
	}{
		{
			name: "OK_1",
			args: args{
				input: domain.Task{
					Title:       "Success",
					Description: "Success test",
					Category:    "Success",
					Complexity:  "Hard",
					Points:      1000,
					Hint:        "Success",
					Flag:        "flag{success}",
					IsActive:    true,
					IsDisabled:  false,
					Author: domain.Author{
						Name:    "SuccessAuthor",
						Contact: "@success",
					},
					Link: "success",
				},
			},
			wantResponse: 1,
			wantErr:      nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			actualResponse, actualErr := repo.SaveTask(testCase.args.input)

			require.Equal(testCase.wantResponse, actualResponse, t)
			require.Equal(testCase.wantErr, actualErr, t)

		})
	}

}

func TestTaskRepository_GetTasks(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		input domain.TaskFilter
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse []*domain.Task
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				input: domain.TaskFilter{
					Id: 1,
				},
			},
			wantResponse: []*domain.Task{{
				Id:          1,
				Title:       "Task1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    false,
				IsDisabled:  true,
				Author: domain.Author{
					Id:      1,
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			},
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				input: domain.TaskFilter{
					IsActive:   true,
					IsDisabled: false,
				},
			},
			wantResponse: []*domain.Task{{
				Id:          2,
				Title:       "Task2",
				Description: "Task2 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 2",
				Flag:        "flag{flag_for_task_2}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Id:      1,
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			}},
			wantErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Поменять если ко-во тасков втесте надо увеличить
			helpCreateTask(t, db, *testCase.wantResponse[0])

			actualResponse, actualErr := repo.GetTasks(testCase.args.input)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}

}

func TestTaskRepository_GetTaskById(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		taskId int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.Task
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				taskId: 1,
			},
			wantResponse: &domain.Task{
				Id:          1,
				Title:       "Task1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Id:      1,
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			helpCreateTask(t, db, *testCase.wantResponse)

			actualResponse, actualErr := repo.GetTaskById(testCase.args.taskId)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		taskId int
		input  domain.TaskUpdate
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse *domain.Task
		wantErr      error
	}{
		// TODO: Add test cases. Обязательно добавить больше тк!!!!!!!!!
		{
			name: "OK",
			args: args{
				taskId: 1,
				input: domain.TaskUpdate{
					Title: "UpdatedTitleTask1",
				},
			},
			wantResponse: &domain.Task{
				Id:          1,
				Title:       "UpdatedTitleTask1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Id:      1,
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			},
			wantErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			task := domain.Task{
				Title:       "Task1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			}
			helpCreateTask(t, db, task)

			actualErr := repo.UpdateTask(testCase.args.taskId, testCase.args.input)
			actualResponse, err := repo.GetTaskById(testCase.args.taskId)

			require.Equal(nil, err, t)
			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTaskRepository_SolveTask(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		taskId int
		teamId int
	}
	testTable := []struct {
		name    string
		args    args
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				taskId: 1,
				teamId: 1,
			},
			wantErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			user := &domain.User{
				Username: "User1",
				Password: "User1Pswd",
				Contact:  "@contact1",
			}
			helpCreateUser(t, db, user)

			team := &domain.Team{
				Name:         "Team1",
				TeamLeaderId: 1,
			}
			helpCreateTeam(t, db, team)

			task := &domain.Task{
				Title:       "Task1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			}
			helpCreateTask(t, db, *task)

			actualErr := repo.SolveTask(testCase.args.taskId, testCase.args.teamId)

			require.Equal(testCase.wantErr, actualErr, 1)
		})
	}
}

func TestTaskRepository_GetTasksSolvedByTeam(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	type args struct {
		teamId           int
		taskIdForSolving int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse []*domain.SolvedTask
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				teamId:           1,
				taskIdForSolving: 1,
			},
			wantResponse: []*domain.SolvedTask{{
				TaskId: 1,
				TeamId: 1,
			}},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			user := &domain.User{
				Username: "User1",
				Password: "User1Pswd",
				Contact:  "@contact1",
			}
			helpCreateUser(t, db, user)

			team := &domain.Team{
				Name:         "Team1",
				TeamLeaderId: 1,
			}
			helpCreateTeam(t, db, team)

			task := &domain.Task{
				Title:       "Task1",
				Description: "Task1 description",
				Category:    "test",
				Complexity:  "hard",
				Points:      1337,
				Hint:        "Hint for task 1",
				Flag:        "flag{flag_for_task_1}",
				IsActive:    true,
				IsDisabled:  false,
				Author: domain.Author{
					Name:    "Author1",
					Contact: "@author1",
				},
				Link: "",
			}
			helpCreateTask(t, db, *task)

			err := repo.SolveTask(testCase.args.taskIdForSolving, testCase.args.teamId)
			actualResponses, actualErr := repo.GetTasksSolvedByTeam(testCase.args.teamId)

			require.Equal(nil, err, t)
			require.Equal(testCase.wantErr, actualErr, t)
			for i, actualResponse := range actualResponses {
				require.Equal(testCase.wantResponse[i].TaskId, actualResponse.TaskId, t)
				require.Equal(testCase.wantResponse[i].TeamId, actualResponse.TeamId, t)
			}
		})
	}
}

func TestTaskRepository_CheckSolvedTask(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	user1 := &domain.User{
		Username: "User1",
		Password: "User1Pswd",
		Contact:  "@contact1",
	}
	helpCreateUser(t, db, user1)

	user2 := &domain.User{
		Username: "User2",
		Password: "User2Pswd",
		Contact:  "@contact2",
	}
	helpCreateUser(t, db, user2)

	team := &domain.Team{
		Name:         "Team1",
		TeamLeaderId: 1,
	}
	helpCreateTeam(t, db, team)

	team2 := &domain.Team{
		Name:         "Team2",
		TeamLeaderId: 2,
	}
	helpCreateTeam(t, db, team2)

	task := &domain.Task{
		Title:       "Task1",
		Description: "Task1 description",
		Category:    "test",
		Complexity:  "hard",
		Points:      1337,
		Hint:        "Hint for task 1",
		Flag:        "flag{flag_for_task_1}",
		IsActive:    true,
		IsDisabled:  false,
		Author: domain.Author{
			Name:    "Author1",
			Contact: "@author1",
		},
		Link: "",
	}
	helpCreateTask(t, db, *task)

	type args struct {
		taskId int
		teamId int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse bool
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				taskId: 1,
				teamId: 2,
			},
			wantResponse: true,
			wantErr:      nil,
		},
		{
			name: "OK_2",
			args: args{
				taskId: 1,
				teamId: 3,
			},
			wantResponse: false,
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.name == "OK_1" {
				err := repo.SolveTask(testCase.args.taskId, testCase.args.teamId)
				require.Equal(nil, err, t)
			}

			actualResponses, actualErr := repo.CheckSolvedTask(testCase.args.taskId, testCase.args.teamId)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponses, t)
		})
	}
}

func TestTaskRepository_SaveTaskSubmission(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	user1 := &domain.User{
		Username: "User1",
		Password: "User1Pswd",
		Contact:  "@contact1",
	}
	helpCreateUser(t, db, user1)

	team := &domain.Team{
		Name:         "Team1",
		TeamLeaderId: 1,
	}
	helpCreateTeam(t, db, team)

	task := &domain.Task{
		Title:       "Task1",
		Description: "Task1 description",
		Category:    "test",
		Complexity:  "hard",
		Points:      1337,
		Hint:        "Hint for task 1",
		Flag:        "flag{flag_for_task_1}",
		IsActive:    true,
		IsDisabled:  false,
		Author: domain.Author{
			Name:    "Author1",
			Contact: "@author1",
		},
		Link: "",
	}
	helpCreateTask(t, db, *task)

	type args struct {
		submission domain.TaskSubmission
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
				submission: domain.TaskSubmission{
					TaskId:    1,
					TeamId:    1,
					UserId:    1,
					Flag:      "flag{HeheBoY}",
					IsCorrect: false,
				},
			},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				submission: domain.TaskSubmission{
					TaskId:    1,
					TeamId:    1,
					UserId:    1,
					Flag:      "flag{falg_for_task_1}",
					IsCorrect: true,
				},
			},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualErr := repo.SaveTaskSubmission(testCase.args.submission)

			require.Equal(testCase.wantErr, actualErr, t)
		})
	}
}

func TestTaskRepository_GetTaskSubmissions(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

	user1 := &domain.User{
		Username: "User1",
		Password: "User1Pswd",
		Contact:  "@contact1",
	}
	helpCreateUser(t, db, user1)

	team := &domain.Team{
		Name:         "Team1",
		TeamLeaderId: 1,
	}
	helpCreateTeam(t, db, team)

	task := &domain.Task{
		Title:       "Task1",
		Description: "Task1 description",
		Category:    "test",
		Complexity:  "hard",
		Points:      1337,
		Hint:        "Hint for task 1",
		Flag:        "flag{flag_for_task_1}",
		IsActive:    true,
		IsDisabled:  false,
		Author: domain.Author{
			Name:    "Author1",
			Contact: "@author1",
		},
		Link: "",
	}
	helpCreateTask(t, db, *task)

	type args struct {
		taskId int
		teamId int
	}
	testTable := []struct {
		name         string
		args         args
		wantResponse []*domain.TaskSubmission
		wantErr      error
	}{
		// TODO: Add test cases.
		{
			name: "OK_1",
			args: args{
				taskId: 1,
				teamId: 2,
			},
			wantResponse: []*domain.TaskSubmission{
				{
					TaskId:    1,
					TeamId:    1,
					UserId:    1,
					Flag:      "flag{HeheBoY}",
					IsCorrect: false,
				},
				{
					TaskId:    1,
					TeamId:    1,
					UserId:    1,
					Flag:      "flag{falg_for_task_1}",
					IsCorrect: true,
				}},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				taskId: 2,
				teamId: 2,
			},
			wantResponse: []*domain.TaskSubmission{},
			wantErr:      nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			TestTaskRepository_SaveTaskSubmission(t)

			actualResponses, actualErr := repo.GetTaskSubmissions(testCase.args.taskId, testCase.args.teamId)

			require.Equal(testCase.wantErr, actualErr, t)
			for i, actualResponse := range actualResponses {
				require.Equal(testCase.wantResponse[i].TaskId, actualResponse.TaskId, t)
				require.Equal(testCase.wantResponse[i].TeamId, actualResponse.TeamId, t)
				require.Equal(testCase.wantResponse[i].Flag, actualResponse.Flag, t)
				require.Equal(testCase.wantResponse[i].IsCorrect, actualResponse.IsCorrect)
				require.Equal(testCase.wantResponse[i].UserId, actualResponse.UserId, t)
			}
		})
	}
}

func helpCreateTask(t *testing.T, db *sqlx.DB, task domain.Task) int {
	t.Helper()

	id, err := postgres.NewTaskRepository(db).SaveTask(task)
	if err != nil {
		t.Fatal(err)
	}

	return id
}
