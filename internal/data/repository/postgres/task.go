package postgres

import (
	"borda/internal/core"
	"borda/internal/core/entity"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("No records found")

type TaskRepository struct {
	db                   *sqlx.DB
	taskTableName        string
	authorTableName      string
	solvedTasksTableName string
}

var _ core.TaskRepository = (*TaskRepository)(nil)

func NewTaskRepository(db *sqlx.DB) core.TaskRepository {
	return TaskRepository{
		db:                   db,
		taskTableName:        "task",
		authorTableName:      "author",
		solvedTasksTableName: "solved_task",
	}
}

func (r TaskRepository) Create(task entity.Task) (int, error) {
	// Create a helper function for preparing failure results.
	fail := func(err error) (int, error) {
		return -1, fmt.Errorf("Create: %v", err)
	}
	// begin transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return fail(err)
	}
	defer tx.Rollback() //nolint

	queryAuthor := fmt.Sprintf("SELECT id FROM public.%s WHERE name = $1 AND contact = $2 LIMIT 1", r.authorTableName)

	var authorId int
	if err := r.db.QueryRowx(queryAuthor, task.Author.Name, task.Author.Contact).Scan(&authorId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			query := fmt.Sprintf("INSERT INTO public.%s (name, contact) VALUES ($1, $2) RETURNING id", r.authorTableName)

			resultRow := tx.QueryRowx(query, task.Author.Name, task.Author.Contact)
			if err := resultRow.Scan(&authorId); err != nil {
				return fail(fmt.Errorf("insert author: %w", err))
			}
		} else {
			return fail(fmt.Errorf("select author: %w", err))
		}
	}

	queryTask := fmt.Sprintf(`
		INSERT INTO public.%s
		(title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`, r.taskTableName)

	var taskId int
	if err = tx.QueryRowx(queryTask, task.Title, task.Description, task.Category,
		task.Complexity, task.Points, task.Hint, task.Flag, task.IsActive,
		task.IsDisabled, authorId).Scan(&taskId); err != nil {
		return fail(fmt.Errorf("insert task %w", err))
	}

	if err := tx.Commit(); err != nil {
		return fail(err)
	}

	return taskId, nil
}

// Get returns one task
func (r TaskRepository) Get(taskId int) (entity.Task, error) {
	query := fmt.Sprintf(`
		SELECT t.id, t.title, t.description, t.category, t.complexity, t.points,
			   t.hint, t.flag, t.is_active, t.is_disabled, a.id AS "author.id",
			   a.name AS "author.name", a.contact AS "author.contact"
		FROM public.%s t, public.%s a
		WHERE t.id = $1 AND t.author_id = a.id
		LIMIT 1`,
		r.taskTableName, r.authorTableName)

	var task entity.Task
	if err := r.db.Get(&task, query, taskId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, fmt.Errorf("no such task: %v", err)
		}
		return entity.Task{}, err
	}
	return task, nil
}

// GetMany returns several task
func (r TaskRepository) FindTasks(filter entity.TaskFilter) ([]entity.Task, error) {
	var filterMap map[string]interface{} = make(map[string]interface{})

	filterJson, err := json.Marshal(filter)
	if err != nil {
		return []entity.Task{}, fmt.Errorf("wrong filter fields: %w", err)
	}

	if err := json.Unmarshal(filterJson, &filterMap); err != nil {
		return []entity.Task{}, fmt.Errorf("wrong filter fields: %w", err)
	}

	params, args := []string{"1 = 1"}, make([]interface{}, 0)
	var index int = 1
	// iterate through filter and build query
	for field, value := range filterMap {
		if field != "offset" && field != "limit" {
			params = append(params, fmt.Sprintf("t.%s = $%d", field, index))
			args = append(args, value)
			index++
		}
	}

	query := fmt.Sprintf(`
		SELECT t.id, t.title, t.description, t.category, t.complexity, t.points,
			   t.hint, t.flag, t.is_active, t.is_disabled, a.id AS "author.id",
			   a.name AS "author.name", a.contact AS "author.contact"
		FROM public.%s t, public.%s a
		WHERE %s AND t.author_id = a.id`,
		r.taskTableName, r.authorTableName, strings.Join(params, " AND "))

	var tasks []entity.Task = make([]entity.Task, 0)

	if err := r.db.Select(&tasks, query, args...); err != nil {
		return tasks, fmt.Errorf("failed select tasks: %w", err)
	}

	if len(tasks) <= 0 {
		return tasks, ErrNotFound
	}

	return tasks, nil
}

func (r TaskRepository) Update(taskId int, newTask entity.Task) error {
	var taskMap map[string]interface{} = make(map[string]interface{})

	jsonTask, err := json.Marshal(newTask)
	if err != nil {
		return fmt.Errorf("bad fields: %w", err)
	}

	if err := json.Unmarshal(jsonTask, &taskMap); err != nil {
		return fmt.Errorf("wrong filter fields: %w", err)
	}

	params, args := []string{"1 = 1"}, make([]interface{}, 0)
	var index int = 1
	// iterate through filter and build query
	for field, value := range taskMap {
		if field != "id" && field != "author_id" && field != "author_name" && field != "author_contact" {
			params = append(params, fmt.Sprintf("t.%s = $%d", field, index))
			args = append(args, value)
			index++
		}
	}

	//TODO: Update Author if author with name and contact doesn't exist.

	query := fmt.Sprintf("UPDATE public.%s SET %s WHERE id = $%d", r.taskTableName, strings.Join(params, ", "), index)
	fmt.Println(query)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query, args, taskId); err != nil {
		return err
	}

	return tx.Commit()
}

// Solve creates a record that the team solved the task
func (r TaskRepository) Solve(taskId, teamId int) error {
	query := fmt.Sprintf("INSERT INTO public.%s (task_id, team_id) VALUES ($1, $2)", r.solvedTasksTableName)

	if _, err := r.db.Exec(query, taskId, teamId); err != nil {
		return err
	}

	return nil
}
