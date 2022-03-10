package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks")
	tasks.Get("", h.getAllTasks)
	tasks.Post("", h.createNewTask)

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask)

	task.Post("/submissions", h.createNewSubmission)
	task.Get("/submissions", h.getAllSubmissions)
}

func (h *Handler) getAllTasks(c *fiber.Ctx) error {
	c.Body()
	return c.JSON(fiber.Map{
		"error": false,
		"tasks": []string{"task1, task2"},
	})
}

func (h *Handler) createNewTask(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Task created",
	})
}

func (h *Handler) updateTask(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}

func (h *Handler) createNewSubmission(c *fiber.Ctx) error {
	return c.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Flag submitted",
	})
}

func (h *Handler) getAllSubmissions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":       false,
		"submissions": []string{"s1", "s2", "s3"},
	})
}
