package api

import (
	"github.com/gofiber/fiber/v2"
)

func initHealthCheckRoute(router fiber.Router) {
	ping := router.Group("/ping")
	ping.Get("/", pongFunction)
}

func pongFunction(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "pong",
	})
}
