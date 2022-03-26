package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAdminRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", h.adminPermissionRequired)

	tasks.Get("", h.getAllTasksAdmin)
	tasks.Post("", h.createNewTask)

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask)
}

func (h *Handler) getAllTasksAdmin(ctx *fiber.Ctx) error {
	var filter domain.TaskFilter
	var tasks []*domain.Task

	err := ctx.BodyParser(&filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	tasks, err = h.AdminService.ShowAllTasks(filter)
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

func (h *Handler) updateTask(ctx *fiber.Ctx) error {

	taskId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusConflict, err.Error())
	}

	var taskData domain.TaskUpdate

	err = ctx.BodyParser(&taskData)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	//TODO: VALIDATE taskData

	err = h.AdminService.UpdateTask(taskId, taskData)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}

func (h *Handler) createNewTask(ctx *fiber.Ctx) error {
	var task domain.Task

	err := ctx.BodyParser(&task)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	//TODO: VALIDATE task

	var taskId int

	taskId, err = h.AdminService.CreateNewTask(task)
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
