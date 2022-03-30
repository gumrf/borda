package domain

import (
	"encoding/json"
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
	AuthorName    string `json:"-"`
	AuthorContact string `json:"-"`
}

// TaskFilter represents a filter passed to FindTasks().
type TaskFilter struct {
	Id         int    `json:"id,omitempty"`
	Category   string `json:"category,omitempty"`
	Complexity string `json:"complexity,omitempty"`
	Points     string `json:"points,omitempty"`
	IsActive   bool   `json:"is_active,omitempty"`
	IsDisabled bool   `json:"is_disabled,omitempty"`
	AuthorId   int    `json:"author_id,omitempty"`

	Offset int `json:"-"`
	Limit  int `json:"-"`
}

//type SolvedTask struct {
//	TaskId    int       `json:"taskId"`
//	TeamId    int       `json:"teamId"`
//	Timestamp time.Time `json:"timestamp"`
//}
//
//type SolvedTasks []SolvedTask

type TaskSubmission struct {
	TaskId    int       `json:"taskId" db:"task_id"`
	TeamId    int       `json:"teamId" db:"team_id"`
	UserId    int       `json:"userId" db:"user_id"`
	Flag      string    `json:"flag" db:"flag"`
	IsCorrect bool      `json:"isCorrect" db:"is_correct"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

// UserTaskResponse содержит только то, что должен увидеть пользователь решающий таски
type UserTaskResponse struct {
	Id          int                      `json:"id"`
	Title       string                   `json:"title" `
	Description string                   `json:"description" `
	Category    string                   `json:"category" `
	Complexity  string                   `json:"complexity" `
	Points      int                      `json:"points" `
	Hint        string                   `json:"hint,omitempty"`
	IsSolved    bool                     `json:"isSolved"`
	Submissions []TaskSubmissionResponse `json:"submissions"`
	Author      Author                   `json:"author"`
}

type TaskSubmissionResponse struct {
	Username  string    `json:"username"`
	Flag      string    `json:"flag" `
	IsCorrect bool      `json:"isCorrect"`
	Timestamp time.Time `json:"timestamp"`
}

type SubmitFlagRequest struct {
	//TaskId int `json:"taskId"`
	//TeamId int    `json:"teamId"`
	//UserId int    `json:"userId"`
	Flag string `json:"flag"`
}

//type TaskSubmissions []TaskSubmission

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
