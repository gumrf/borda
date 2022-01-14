package entity

type Role struct {
	Id   int    `json:"roleId"`
	Name string `json:"roleName"`
}

type UserRoles struct {
	UserId int `json:"userId"`
	RoleId int `json:"roleId"`
}
