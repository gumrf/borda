package core

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"name"`
	Password string `json:"password" db:"password"`
	Contact  string `json:"contact" db:"contact"`
	TeamId   int    `json:"teamId" db:"team_id"`
}

type UserRepository interface {
	CrudRepository[User]
	// GetUserByCredentials ...
	FindByCredentials(username, password string) (User, error)
	// UpdatePassword Change user password
	UpdatePassword(userId int, newPassword string) error
	// GrantRole assigns a role to a user
	GrantRole(userId, roleId int) error
}
