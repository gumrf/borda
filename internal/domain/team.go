package domain

import "github.com/google/uuid"

type Team struct {
	Id          int       `json:"teamId"`
	Name        string    `json:"teamName"`
	TeamLeader  User      `json:"teamLeader"`
	Token       uuid.UUID `json:"token"`
	TeamMembers []User    `json:"teamMembers"`
}
