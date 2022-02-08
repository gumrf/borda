package postgres_test

import (
	"borda/internal/core/entity"
	"borda/internal/data/repository/postgres"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestTaskRepository(t *testing.T) {
// 	t.Run("Create", func(t *testing.T) { testTaskRepository_Create(t, r) })
// 	t.Run("Get", func(t *testing.T) { testTaskRepository_Get(t, r) })
// 	t.Run("GetMany", func(t *testing.T) { testTaskRepository_GetMany(t, r) })
// 	t.Run("Update", func(t *testing.T) { testTaskRepository_Update(t, r) })
// }

func Test_TaskRepository_Create(t *testing.T) {
	var task = entity.Task{
		Id:          1,
		Title:       "Easy tcp",
		Description: "lorem bla bla bla",
		Category:    "web",
		Complexity:  "easy",
		Points:      100,
		Hint:        "Keep searching",
		Flag:        "flag{flag_is_here}",
		IsActive:    true,
		IsDisabled:  false,
		Author: entity.Author{
			Id:      1,
			Name:    "Max",
			Contact: "max@gmail.com",
		},
	}

	db := MustConnectAndMigrate(t)
	repo := postgres.NewTaskRepository(db)

	taskId, err := repo.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	query := `
		SELECT t.id, t.title, t.description, t.category, t.complexity, t.points,
			   t.hint, t.flag, t.is_active, t.is_disabled, a.id AS "author.id",
			   a.name AS "author.name", a.contact AS "author.contact"
		FROM public.task t, public.author a
		WHERE t.id = $1 AND t.author_id = a.id
		LIMIT 1`

	var testTask entity.Task
	if err := db.Get(&testTask, query, taskId); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, task, testTask, fmt.Sprintf("Got: %v\nWant: %v", testTask, task))
}

func Test_TaskRepository_Get(t *testing.T) {
	var task = entity.Task{
		Id:          1,
		Title:       "Easy tcp",
		Description: "lorem bla bla bla",
		Category:    "web",
		Complexity:  "easy",
		Points:      100,
		Hint:        "Keep searching",
		Flag:        "flag{flag_is_here}",
		IsActive:    true,
		IsDisabled:  false,
		Author: entity.Author{
			Id:      1,
			Name:    "Max",
			Contact: "max@gmail.com",
		},
	}

	db := MustConnectAndMigrate(t)
	repo := postgres.NewTaskRepository(db)

	taskId, err := repo.Create(task)
	if err != nil {
		t.Fatal("Oops..Can't create task: ", err)
	}

	testTask, err := repo.Get(taskId)
	if err != nil {
		t.Fatal("Oops..Can't get task: ", err)
	}

	t.Logf("Is equal: %v\n", reflect.DeepEqual(task, testTask))
	assert.Equal(t, task, testTask, fmt.Sprintf("Want: %v\nGot: %v\n", task, testTask))
}

func Test_TaskRepository_FindTask(t *testing.T) {
	db := MustConnectAndMigrate(t)
	repo := postgres.NewTaskRepository(db)

	var tasks []entity.Task
	for i := 0; i < 10; i++ {
		var task = entity.Task{
			Id:          1 + i,
			Title:       "Easy tcp",
			Description: "lorem bla bla bla",
			Category:    "web",
			Complexity:  "easy",
			Points:      100 + (i+1%2)*100,
			Hint:        "Keep searching",
			Flag:        "flag{flag_is_here}",
			IsActive:    true,
			IsDisabled:  false,
			Author: entity.Author{
				Id:      1,
				Name:    "Max",
				Contact: "max@gmail.com",
			},
		}

		tasks = append(tasks, task)

		_, err := repo.Create(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("ErrNotFound", func(t *testing.T) {
		_, err := repo.FindTasks(entity.TaskFilter{
			Id:         5,
			Category:   "web",
			Complexity: "easy",
			Points:     "200",
			IsActive:   true,
			IsDisabled: false,
			Offset:     0,
			Limit:      0,
		})

		assert.ErrorIs(t, postgres.ErrNotFound, err)
	})

	t.Run("FindOne", func(t *testing.T) {
		testTasks, err := repo.FindTasks(entity.TaskFilter{
			Id:         6,
			Category:   "web",
			Complexity: "easy",
			Points:     "700",
			IsActive:   true,
			IsDisabled: false,
			Offset:     0,
			Limit:      0,
		})

		assert.NoError(t, err, "should be not nil")
		assert.Equal(t, []entity.Task{tasks[5]}, testTasks)
	})

	t.Run("FindByCategory", func(t *testing.T) {
		testTasks, err := repo.FindTasks(entity.TaskFilter{
			Category: "web",
		})

		assert.NoError(t, err, "should be not nil")
		assert.Equal(t, tasks, testTasks)
	})

}
