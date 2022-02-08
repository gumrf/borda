package entity

import (
	"time"
)

// Task is task
type Task struct {
	Id          int    `json:"id" db:"id"`
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
}

// TaskFilter represents a filter passed to FindTasks().
type TaskFilter struct {
	// Filtering fields.
	Id         int    `json:"id,omitempty"`
	Category   string `json:"category,omitempty"`
	Complexity string `json:"complexity,omitempty"`
	Points     string `json:"points,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	IsDisabled bool   `json:"is_disabled,omitempty"`
	AuthorId   int    `json:"author_id,omitempty"`

	// Restrict to subset of results.
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type Author struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Contact string `json:"contact" db:"contact"`
}

type SolvedTask struct {
	TaskId    int       `json:"taskId"`
	TeamId    int       `json:"teamId"`
	Timestamp time.Time `json:"timestemp"`
}

type SolvedTasks []SolvedTask

type TaskSubmission struct {
	TaskId    int       `json:"taskId"`
	TeamId    int       `json:"teamId"`
	UserId    int       `json:"userId"`
	Flag      string    `json:"flag"`
	IsCorrect bool      `json:"isCorrect"`
	Timestemp time.Time `json:"timestemp"`
}

type TaskSubmissions []TaskSubmission
