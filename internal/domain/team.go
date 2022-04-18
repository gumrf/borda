package domain

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

type AttachTeamInput struct {
	// Method can be 'create' or 'join'
	Method string `json:"method"`
	// Attribute depends on method.
	// Can be team token or team name.
	Attribute string `json:"attribute"`
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
