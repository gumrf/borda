package core

import (
	"encoding/json"
	"time"
)

// TODO: Refactoring

type TaskRepository interface {
	SaveTask(task Task) (int, error)
	GetTaskById(id int) (*Task, error)
	GetTasks(TaskFilter) ([]*Task, error)
	GetTasksSolvedByTeam(teamId int) ([]*SolvedTask, error)
	UpdateTask(id int, update TaskUpdate) error
	SolveTask(taskId, teamId int) error
	SaveTaskSubmission(submission TaskSubmission) error
	GetTaskSubmissions(taskId, teamId int) ([]*TaskSubmission, error)
	CheckSolvedTask(taskId, teamId int) (bool, error)
}

type Task struct {
	Id          int    `json:"id,omitempty" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Category    string `json:"category" db:"category"`
	Complexity  string `json:"complexity" db:"complexity"`
	Points      int    `json:"points" db:"points"`
	Hint        string `json:"hint" db:"hint"`
	Flag        string `json:"flag" db:"flag"`
	IsActive    bool   `json:"isActive" db:"is_active"`
	IsDisabled  bool   `json:"isDisabled" db:"is_disabled"`
	Author      Author `json:"author" db:"author"`
	Link        string `json:"link" db:"link"`
}

type Author struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Contact string `json:"contact" db:"contact"`
}

type TaskUpdate struct {
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Category      string `json:"category,omitempty"`
	Complexity    string `json:"complexity,omitempty"`
	Points        int    `json:"points,omitempty"`
	Hint          string `json:"hint,omitempty"`
	Flag          string `json:"flag,omitempty"`
	IsActive      string `json:"is_active,omitempty"`   // Обязательно is_active
	IsDisabled    string `json:"is_disabled,omitempty"` // Обязательно is_disabled
	AuthorName    string `json:"-"`
	AuthorContact string `json:"-"`
	Link          string `json:"link,omitempty"`
}

type TaskFilter struct {
	Id         int    `json:"id,omitempty"`
	Category   string `json:"category,omitempty"`
	Complexity string `json:"complexity,omitempty"`
	Points     string `json:"points,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`   // Обязательно is_active
	IsDisabled bool   `json:"is_disabled,omitempty"` // Обязательно is_disabled
	AuthorId   int    `json:"author_id,omitempty"`

	Offset int `json:"-"`
	Limit  int `json:"-"`
}

type TaskSubmission struct {
	TaskId    int       `json:"taskId" db:"task_id"`
	TeamId    int       `json:"teamId" db:"team_id"`
	UserId    int       `json:"userId" db:"user_id"`
	Flag      string    `json:"flag" db:"flag"`
	IsCorrect bool      `json:"isCorrect" db:"is_correct"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type SolvedTask struct {
	TaskId    int       `json:"taskId" db:"task_id"`
	TeamId    int       `json:"teamId" db:"team_id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func (u *TaskUpdate) ToMap() (map[string]interface{}, error) {
	bytes, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *TaskFilter) ToMap() (map[string]interface{}, error) {
	bytes, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
