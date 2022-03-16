package api

import (
	"borda/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", AuthRequired)
	tasks.Get("", h.getAllTasks)    // ++
	tasks.Post("", h.createNewTask) // ++

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask) //Only admin изменение таска - вариотивная приходит структура см. домен

	task.Post("/submissions", h.createNewSubmission) // Users Приходит флаг
	task.Get("/submissions", h.getAllSubmissions)    // Users Вернуть все попытки решения
}

func (h *Handler) getAllTasks(ctx *fiber.Ctx) error {
	var filter domain.TaskFilter
	var tasks []*domain.Task

	err := ctx.BodyParser(&filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	tasks, err = h.UserUsecase.ShowAllTasks(filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	type TaskRespose struct {
		Response []*domain.Task
	}

	return ctx.Status(fiber.StatusOK).JSON(TaskRespose{
		Response: tasks,
	})
}

func (h *Handler) createNewTask(ctx *fiber.Ctx) error {
	var task domain.Task

	err := ctx.BodyParser(&task)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	// VALIDATE task

	var taskId int

	taskId, err = h.AdminUsecase.CreateNewTask(task)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": "Task created",
		"task_id": taskId,
	})
}

func (h *Handler) updateTask(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}

func (h *Handler) createNewSubmission(ctx *fiber.Ctx) error {
	/*
		Приходит флаг, таск, тима, юзер
	*/
	var submission domain.SubmitTaskRequest
	err := ctx.BodyParser(&submission)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	//VALIDATE submission.Flag

	var message string
	message, err = h.UserUsecase.TryToSolveTask(submission)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusConflict, err.Error())
	}

	return ctx.Status(201).JSON(fiber.Map{
		"error":   false,
		"message": message,
		"flag":    submission.Flag,
	})
}

func (h *Handler) getAllSubmissions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":       false,
		"submissions": []string{"s1", "s2", "s3"},
	})
}
