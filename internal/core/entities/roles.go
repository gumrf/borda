package entities

type Role struct {
	ID   int    `json:"roleId"`
	Name string `json:"roleName"`
}

type UserRole struct {
	UserID int
	User   User

	RoleID int
	Role   Role
}
