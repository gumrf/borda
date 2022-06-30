package team

import "borda/internal/pkg/core"

type Team = core.Team
type User = core.User

type TeamMember struct {
	UserId int    `json:"id" db:"user_id"`
	Name   string `json:"name" db:"user_name"`
}

type TeamResponse struct {
	Id      int          `json:"id"`
	Name    string       `json:"name"`
	Token   string       `json:"token"`
	Captain TeamMember   `json:"captain"`
	Members []TeamMember `json:"members"`
}

type Scoreboard struct {
	TeamName         string `json:"teamName"`
	Score            int    `json:"score"`
	TeamMembersCount int    `json:"teamMembersCount"`
}
