package api

import (
	"github.com/gofiber/fiber/v2"
)

func initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks")
	tasks.Get("", getAllTasks)
	tasks.Post("/:id", updateTask)
}

func getAllTasks(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
		"tasks": []string{"task1, task2"},
	})
}

func updateTask(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}
