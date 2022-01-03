package domain

import (
	"time"
)

type Task struct {
	Id         int      `json:"id"`
	Title      string   `json:"title"`
	Decription string   `json:"description"`
	Category   string   `json:"category"`
	Complexity string   `json:"complexity"`
	Pionts     int      `json:"points"`
	Hint       string   `json:"hint"`
	Flag       string   `json:"flag"`
	IsActive   bool     `json:"isActive"`
	IsDisabled bool     `json:"isDisabled"`
	Authors    []Author `json:"authors"`
}

type Author struct {
	Id      int    `json:"authorId"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type SolvedTask struct {
	TaskId    int       `json:"taskId"`
	TeamId    int       `json:"teamId"`
	Timestamp time.Time `json:"timestemp"`
}

type SolvedTasks []SolvedTask

type TaskSubmission struct {
	TaskId        int       `json:"taskId"`
	TeamId        int       `json:"teamId"`
	SubmitionerId int       `json:"submitionerId"`
	Submission    string    `json:"submission"`
	IsCorrect     bool      `json:"isCorrect"`
	Timestemp     time.Time `json:"timestemp"`
}

type TaskSubmissions []TaskSubmission
