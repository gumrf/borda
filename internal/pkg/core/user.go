package core

import "context"

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
	FindByCredentials(ctx context.Context, username, password string) (User, error)
	// UpdatePassword Change user password
	UpdatePassword(ctx context.Context, userId int, newPassword string) (User, error)
	// GrantRole assigns a role to a user
	GrantRole(ctx context.Context, userId, roleId int) error
}
