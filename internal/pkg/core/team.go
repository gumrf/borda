package core

import (
	"context"
)

type TeamRepository interface {
	CrudRepository[Team]
	FindByToken(token string) (Team, error)
	FindByUserId(userId int) (Team, error)
	AddMember(ctx context.Context, teamId, userId int) error
	IsTeamFull(ctx context.Context, teamId int) bool
}

type Team struct {
	Id           int          `json:"id" db:"id"`
	Name         string       `json:"teamName" db:"name"`
	TeamLeaderId int          `json:"teamLeaderId" db:"team_leader_id"`
	Token        string       `json:"token" db:"token"`
	Members      []TeamMember `json:"members" db:"members"`
}

type TeamMember struct {
	UserId int    `json:"id" db:"user_id"`
	Name   string `json:"name" db:"user_name"`
}
