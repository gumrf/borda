package entity

import (
	"github.com/google/uuid"
)

type Team struct {
	Name         string    `json:"teamName"`
	TeamLeaderId int       `json:"teamLeaderId"`
	Token        uuid.UUID `json:"token"`
	TeamMembers  []User    `json:"teamMembers" gorm:""`
}
