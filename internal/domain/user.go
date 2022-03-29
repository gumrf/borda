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
	TeamName string `json:"team_name,omitempty"`
}
