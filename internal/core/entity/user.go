package entity

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	TeamId   int    `json:"teamId"`
}
