package domain

type Team struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"teamName" db:"name"`
	TeamLeaderId int    `json:"teamLeaderId" db:"team_leader_id"`
	Token        string `json:"token" db:"token"`
}

type AttachTeamInput struct {
	// Method can be 'create' or 'join'
	Method string `json:"method"`
	// Attribute depends on method.
	// Can be team token or team name.
	Attribute string `json:"attribute"`
}
