package user

import "borda/internal/pkg/core"

// type User struct {
// 	Id       int    `json:"id" db:"id"`
// 	Username string `json:"username" db:"name"`
// 	Password string `json:"password" db:"password"`
// 	Contact  string `json:"contact" db:"contact"`
// 	TeamId   int    `json:"teamId" db:"team_id"`
// }

type User = core.User

type PublicUserProfileResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Team     struct {
		Id   int    `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"team"`
}

type PrivateUserProfileResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Contact  string `json:"contact"`
	// Team     TeamResponse `json:"team"`
}
