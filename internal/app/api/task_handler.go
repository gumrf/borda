package api

import (
	"github.com/gofiber/fiber/v2"
)

func initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks")
	tasks.Get("", getAllTasks)
	tasks.Post("", createNewTask)

	task := router.Group("/tasks/:id")
	task.Patch("", updateTask)

	task.Post("/submissions", createNewSubmission)
	task.Get("/submissions", getAllSubmissions)
}

func getAllTasks(c *fiber.Ctx) error {
	c.Body()
	return c.JSON(fiber.Map{
		"error": false,
		"tasks": []string{"task1, task2"},
	})
}

func createNewTask(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Task created",
	})
}

func updateTask(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}

func createNewSubmission(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Flag submitted",
	})
}

func getAllSubmissions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":       false,
		"submissions": []string{"s1", "s2", "s3"},
	})
}
