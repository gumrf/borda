package domain

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"name"`
	Password string `json:"password" db:"password"`
	Contact  string `json:"contact" db:"contact"`
	TeamId   int    `json:"teamId" db:"team_id"`
}

type SignUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	Username string `json:"username"`
	TeamName string `json:"teamName,omitempty"`
}

type UserProfileResponse struct {
	Id          int              `json:"userId"`
	Username    string           `json:"username"`
	Contact     string           `json:"contact"`
	TeamId      int              `json:"teamId"`
	TeamName    string           `json:"teamName"`
	TeamMembers []MemberResponse `json:"teamMembers"`
}
