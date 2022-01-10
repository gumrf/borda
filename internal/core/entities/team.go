package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name        string    `json:"teamName"`
	TeamLeader  User      `json:"teamLeader"`
	Token       uuid.UUID `json:"token"`
	TeamMembers []User    `json:"teamMembers" gorm:"many2many:team_members;"`
}
