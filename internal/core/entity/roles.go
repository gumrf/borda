package entity

type Role struct {
	Id   int    `json:"roleId"`
	Name string `json:"roleName"`
}

type UserRoles []Role
