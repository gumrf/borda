package api

import (
	"github.com/gofiber/fiber/v2"
)

func initUserRoutes(router fiber.Router) {
	users := router.Group("/users")
	users.Get("", getAllUsers)
	users.Get("/:id", getUserById)
}

func getUserById(c *fiber.Ctx) error {
	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":  false,
		"userId": c.Params("id"),
	})
}

func getAllUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"users": []string{"user1, user2"},
	})
}
