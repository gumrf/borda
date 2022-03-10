package api

import (
	"github.com/gofiber/fiber/v2"
)

func initAuthRoute(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/signIn", signIn)
}

func signUp(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"token": "JWT-Token",
	})
}

func signIn(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"token": "JWT-Token",
	})
}

func signOut(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
	})
}
