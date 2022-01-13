package entities

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	TeamId   int    `json:"teamId"`
}
