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

type PublicUserProfileResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Team     struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
}

type PrivateUserProfileResponse struct {
	Id       int          `json:"id"`
	Username string       `json:"username"`
	Contact  string       `json:"contact"`
	Team     TeamResponse `json:"team"`
}

type SignInResponse struct {
	Token string `json:"token"`
}
