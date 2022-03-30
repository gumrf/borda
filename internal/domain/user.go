package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"name"`
	Password string `json:"password" db:"password"`
	Contact  string `json:"contact" db:"contact"`
	TeamId   int    `json:"teamId" db:"team_id"`
}

type UserSignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
}

type UserSignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username string `json:"username"`
	TeamName string `json:"teamName,omitempty"`
}

// Переименовать во что-то нормальное, жду предложения номальных названий этих двух структур :)
type UserProfileResponse struct {
	Username     string                `json:"username"`
	Contact      string                `json:"contact"`
	TeamName     string                `json:"teamName"`
	IsTeamLeader bool                  `json:"teamLeaderName"`
	TeamMembers  []TeamMembersResponse `json:"teamMembers,omitempty"`
}
