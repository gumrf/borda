package postgres_test

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskRepository_SaveTask(t *testing.T) {

	db := MustOpenDB(t)
	repo := postgres.NewTaskRepository(db)
	assert := assert.New(t)

	type args struct {
		input domain.Task
	}

	testTable := []struct {
		name    string
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name: "OK",
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
						Id:      1,
						Name:    "Author1",
						Contact: "@author1",
					},
					Link: "success",
				},
			},
			wantId: 4,
		},
		{
			name: "fail",
			args: args{
				input: domain.Task{
					//Title:       "",
					Description: "",
					Category:    "",
					Complexity:  "",
					Points:      0,
					Hint:        "",
					Flag:        "",
					IsActive:    false,
					//IsDisabled:  false,
					Author: domain.Author{
						//Id: 1,
					},
					Link: "",
				},
			},
			wantId: 5,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			tx, err := db.Beginx()
			if err != nil {
				t.Errorf("error with trancaction = %v", err)
			}
			defer tx.Rollback()

			got, err := repo.SaveTask(testCase.args.input)
			if testCase.wantErr {
				t.Error(err)
				assert.Error(err, t)
			} else {
				assert.NoError(err, t)
				assert.Equal(testCase.wantId, got, "equal")
			}
		})
	}

}

func TestTaskRepository_findOrCreateAuthor(t *testing.T) {

	db := MustOpenDB(t)
	repo := postgres.NewTaskRepository(db)
	assert := assert.New(t)

	type args struct {
		author domain.Author
	}
	testTable := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: args{
				author: domain.Author{
					Name:    "Tester",
					Contact: "@tester",
				},
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			tx, err := db.Beginx()
			if err != nil {
				t.Errorf("error with trancaction = %v", err)
			}
			defer tx.Rollback()

			got, err := repo.FindOrCreateAuthor(tx, testCase.args.author)
			if testCase.wantErr {
				assert.Error(err, t)
			} else {
				assert.NoError(err, t)
				assert.Equal(testCase.want, got, "equal")
			}
		})
	}
}
