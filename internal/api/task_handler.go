package api

import (
	"borda/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", AuthRequired)
	tasks.Get("", h.getAllTasks)    // Users вернуть массив тасков для отображения
	tasks.Post("", h.createNewTask) //Only admin просто создать новый таск

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask) //Only admin изменение таска - вариотивная приходит структура см. домен

	task.Post("/submissions", h.createNewSubmission) // Users Приходит флаг
	task.Get("/submissions", h.getAllSubmissions)    // Users Вернуть все попытки решения
}

func (h *Handler) getAllTasks(ctx *fiber.Ctx) error {
	ctx.Body()

	var tasks []*domain.Task

	tasks, err := h.UserUsecase.ShowAllTasks()
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	//return ctx.JSON(fiber.Map{
	//	"error": false,
	//	"tasks": []string{"task1, task2"},
	//})
	//INSERT INTO task (title, description, category, complexity, points, hint, flag, is_active, is_disabled, author_id) VALUES ('task1', 'This is description', 'penis', 'hard', 50, 'this is hint', 'Allah', true, true, 1);

	//var output []byte
	//
	//for t := range tasks {
	//	output, err = json.Marshal(t)
	//	if err != nil {
	//		return NewErrorResponse(ctx, fiber.StatusConflict, err.Error())
	//	}
	//}

	return ctx.Status(fiber.StatusOK).JSON(tasks)
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
