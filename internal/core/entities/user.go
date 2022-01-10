package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Contact  string `json:"contact"`
	TeamID   int    `json:"teamId"`
}
