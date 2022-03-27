package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users", h.authRequired, h.checkUserInTeam)
	users.Get("", h.getAllUsers)

	user := users.Group("/:id")
	user.Get("", h.getUserById)
}

func (h *Handler) getUserById(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":  false,
		"userId": c.Params("id"),
	})
}

func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"users": []string{"user1, user2"},
	})
}
