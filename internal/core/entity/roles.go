package entity

type Role struct {
	Id   int    `json:"roleId" db:"id"`
	Name string `json:"roleName" db:"name"`
}

type UserRoles struct {
	UserId int `json:"userId" db:"user_id"`
	RoleId int `json:"roleId" db:"role_id"`
}
