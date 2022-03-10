package domain

type Team struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"teamName" db:"name"`
	TeamLeaderId int    `json:"teamLeaderId" db:"team_leader_id"`
	Token        string `json:"token" db:"token"`
}
