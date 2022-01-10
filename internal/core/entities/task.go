package entities

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	// Id         int      `json:"id"`
	Title      string `json:"title" gorm:"not null"`
	Decription string `json:"description"`
	Category   string `json:"category"`
	Complexity string `json:"complexity"`
	Pionts     int    `json:"points"`
	Hint       string `json:"hint"`
	Flag       string `json:"flag"`
	IsActive   bool   `json:"isActive"`
	IsDisabled bool   `json:"isDisabled"`

	// AuthorID *int     `json:"-"`
	Authors []Author `json:"authors" gorm:"many2many:team_members;"`
}

type Author struct {
	ID      int    `json:"authorId"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type SolvedTask struct {
	TaskID int  `json:"taskId"`
	Task   Task `json:"task"`

	TeamID int  `json:"teamId"`
	Team   Team `json:"team"`

	Timestamp time.Time `json:"timestemp"`
}

type SolvedTasks []SolvedTask

type TaskSubmission struct {
	TaskID int  `json:"taskId"`
	Task   Task `json:"task"`

	TeamID int  `json:"teamId"`
	Team   Team `json:"team"`

	UserId int  `json:"userId"`
	User   User `json:"user"`

	Flag      string    `json:"flag"`
	IsCorrect bool      `json:"isCorrect"`
	Timestemp time.Time `json:"timestemp"`
}

type TaskSubmissions []TaskSubmission
