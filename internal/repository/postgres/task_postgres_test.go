package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTaskRepository_SaveTask(t *testing.T) {
	db := MustOpenDB(t)
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
			wantResponse: 4,
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
				IsActive:    true,
				IsDisabled:  false,
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
			}},
			wantErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualResponse, actualErr := repo.GetTasks(testCase.args.input)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}

}

func TestTaskRepository_GetTaskById(t *testing.T) {
	db := MustOpenDB(t)
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
			actualResponse, actualErr := repo.GetTaskById(testCase.args.taskId)

			require.Equal(testCase.wantErr, actualErr, t)
			require.Equal(testCase.wantResponse, actualResponse, t)
		})
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	db := MustOpenDB(t)
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
		// TODO: Add test cases.
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
				teamId: 2,
			},
			wantErr: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actualErr := repo.SolveTask(testCase.args.taskId, testCase.args.teamId)

			require.Equal(testCase.wantErr, actualErr, 1)
		})
	}
}

func TestTaskRepository_GetTasksSolvedByTeam(t *testing.T) {
	db := MustOpenDB(t)
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
				teamId:           2,
				taskIdForSolving: 1,
			},
			wantResponse: []*domain.SolvedTask{{
				TaskId: 1,
				TeamId: 2,
			}},
			wantErr: nil,
		},
		{
			name: "OK_2",
			args: args{
				teamId:           3,
				taskIdForSolving: 1,
			},
			wantResponse: []*domain.SolvedTask{{
				TaskId: 1,
				TeamId: 3,
			}},
			wantErr: nil,
		},
		{
			name: "OK_3",
			args: args{
				teamId:           3,
				taskIdForSolving: 2,
			},
			wantResponse: []*domain.SolvedTask{{
				TaskId: 1,
				TeamId: 3,
			},
				{
					TaskId: 2,
					TeamId: 3,
				}},
			wantErr: nil,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
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
	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

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
	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

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
					TeamId:    2,
					UserId:    2,
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
					TeamId:    2,
					UserId:    2,
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
	repo := postgres.NewTaskRepository(db)
	require := require.New(t)

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
			wantResponse: []*domain.TaskSubmission{{
				TaskId:    1,
				TeamId:    2,
				UserId:    2,
				Flag:      "flag{HeheBoY}",
				IsCorrect: false,
			},
				{
					TaskId:    1,
					TeamId:    2,
					UserId:    2,
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
