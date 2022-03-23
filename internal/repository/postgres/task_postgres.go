package postgres

import (
	"borda/internal/domain"

	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository struct {
	db                      *sqlx.DB
	taskTableName           string
	authorTableName         string
	solvedTasksTableName    string
	taskSubmissionTableName string
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{
		db:                      db,
		taskTableName:           "task",
		authorTableName:         "author",
		solvedTasksTableName:    "solved_task",
		taskSubmissionTableName: "task_submission",
	}
}

func (r TaskRepository) CreateNewTask(task domain.Task) (int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return -1, fmt.Errorf("TaskRepository.Create: beginx: %w", err)
	}
	defer tx.Rollback() //nolint

	_, err = r.findOrCreateAuthor(tx, &task.Author)
	if err != nil {
		return -1, fmt.Errorf("TaskRepository.Create: findOrCreateAuthor: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			title, 
			description, 
			category, 
			complexity, 
			points, 
			hint, 
			flag, 
			is_active, 
			is_disabled, 
			author_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		r.taskTableName)

	result := tx.QueryRowx(
		query,
		task.Title, task.Description, task.Category, task.Complexity,
		task.Points, task.Hint, task.Flag, task.IsActive, task.IsDisabled,
		task.Author.Id)

	if err := result.Scan(&task.Id); err != nil {
		return -1, fmt.Errorf("TaskRepository.Create: scan: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return -1, fmt.Errorf("TaskRepository.Create: commit: %w", err)
	}

	return task.Id, nil
}

func (r TaskRepository) findOrCreateAuthor(tx *sqlx.Tx, author *domain.Author) (int, error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM public.%s
		WHERE name = $1
		LIMIT 1`,
		r.authorTableName,
	)

	result := tx.QueryRowx(query, author.Name)
	err := result.Scan(&author.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_query := fmt.Sprintf(`
				INSERT INTO public.%s (
					name,
					contact
				)
				VALUES ($1, $2)
				RETURNING id`,
				r.authorTableName)

			_result := tx.QueryRowx(_query, author.Name, author.Contact)
			if err := _result.Scan(&author.Id); err != nil {
				return -1, err
			}
		} else {
			return -1, err
		}
	}

	return author.Id, nil
}

// GetTasks returns a list of tasks based on a filter.
func (r TaskRepository) GetTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // nolint

	tasks, err := r.findTasks(tx, filter)
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

// GetTaskById search for task with specified id
func (r TaskRepository) GetTaskById(taskId int) (*domain.Task, error) {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // nolint

	tasks, err := r.findTasks(tx, domain.TaskFilter{Id: taskId})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("can't find task with id=%d", taskId)
		}
		return nil, err
	}

	return tasks[0], nil
}

// findTasks returns a list of matching tasks.
func (r TaskRepository) findTasks(tx *sqlx.Tx, filter domain.TaskFilter) (_ []*domain.Task, err error) {
	f, err := filter.ToMap()
	if err != nil {
		return nil, err
	}

	index, params, args := 1, []string{"1 = 1"}, make([]interface{}, 0)

	for filter, value := range f {
		params = append(params, fmt.Sprintf("task.%s = $%d", filter, index))
		args = append(args, value)
		index++
	}

	query := fmt.Sprintf(`
		SELECT
			task.id, 
			task.title,
			task.description,
			task.category, 
			task.complexity,  
			task.points,
			task.hint,   
			task.flag, 
			task.is_active, 
			task.is_disabled, 
			author.id AS "author.id",
			author.name AS "author.name",
			author.contact AS "author.contact"
		FROM public.%s
		LEFT JOIN public.%s ON task.author_id = author.id
		WHERE %s
		ORDER BY id 
		`+formatLimitOffset(filter.Limit, filter.Offset),
		r.taskTableName, r.authorTableName, strings.Join(params, " AND "))

	tasks := make([]*domain.Task, 0)
	if err := tx.Select(&tasks, query, args...); err != nil {
		return nil, err
	}

	return tasks, nil
}

// FormatLimitOffset returns a SQL string for a given limit & offset.
// Clauses are only added if limit and/or offset are greater than zero.
func formatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)

	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)

	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}

	return ""
}

func (r TaskRepository) UpdateTask(taskId int, update domain.TaskUpdate) error {
	updateFields, err := update.ToMap()
	if err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err

	}
	defer tx.Rollback() // nolint

	// init sql query params & args
	params, args, index := make([]string, 0), make([]interface{}, 0), 1

	for field, value := range updateFields {
		params = append(params, field+" = $"+strconv.Itoa(index))
		args = append(args, value)
		index++
	}

	query := fmt.Sprintf(`
		UPDATE public.%s
		SET %s
		WHERE id = $%d`,
		r.taskTableName,
		strings.Join(params, ", "),
		index,
	)

	args = append(args, taskId)

	fmt.Println(query)

	if _, err := tx.Exec(query, args...); err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
			UPDATE public.%[1]s
			SET name = $1, contact = $2
			WHERE id = (
				SELECT id
				FROM public.%[1]s
				WHERE name = $1 AND contact = $2
			)`,
		r.authorTableName),
		update.AuthorName,
		update.AuthorContact)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// SolveTask creates a record that the team solved the task
func (r TaskRepository) SolveTask(taskId, teamId int) error {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			task_id,
			team_id
		)
		VALUES ($1, $2)`,
		r.solvedTasksTableName)

	if _, err := r.db.Exec(query, taskId, teamId); err != nil {
		return err
	}

	return nil
}

func (r TaskRepository) ChekSolvedTask(taskId, teamId int) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 FROM public.%s
			WHERE task_id=$1 AND team_id=$2
		)`,
		r.solvedTasksTableName)

	var isTaskSolved bool

	err := r.db.Get(&isTaskSolved, query, taskId, teamId) //Return false if task not solved
	if err != nil {
		return isTaskSolved, err
	}

	return isTaskSolved, nil
}

func (r TaskRepository) FillTaskSubmission(value domain.SubmitTaskRequest, isCorrect bool) error {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			task_id,
			team_id,
			user_id,
			flag,
			is_correct
		)
		VALUES ($1, $2, $3, $4, $5)`,
		r.taskSubmissionTableName)

	if _, err := r.db.Exec(query, value.TaskId, value.TeamId, value.UserId,
		value.Flag, isCorrect); err != nil {
		return err
	}
	return nil
}

func (r TaskRepository) ShowTaskSubmissions(value domain.SubmitTaskRequest) ([]*domain.TaskSubmission, error) {
	ctx := context.Background()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := fmt.Sprintf(`
	SELECT * FROM public.%s 
	WHERE task_submission.team_id=$1 AND task_submission.task_id=$2`,
		r.taskSubmissionTableName)

	result := make([]*domain.TaskSubmission, 0)

	if err := tx.Select(&result, query, value.TeamId, value.TaskId); err != nil {
		return nil, err
	}

	return result, nil
}
