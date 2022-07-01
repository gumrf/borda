package core

import (
	"context"
)

type TeamRepository interface {
	CrudRepository[Team]
	FindByToken(ctx context.Context, token string) (Team, error)
	FindByUserId(ctx context.Context, userId int) (Team, error)
	AddMember(ctx context.Context, teamId, userId int) error
	IsTeamFull(ctx context.Context, teamId int) bool
}

type TeamService interface {
	JoinTeam(ctx context.Context, userId int, token string) error
	CreateTeam(ctx context.Context, userId int, name string) error
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
	// IsCaptain bool   `db:"is_captain"`
}
