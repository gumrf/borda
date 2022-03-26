package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAdminRoutes(router fiber.Router) {
	admin := router.Group("/admin", h.adminPermissionRequired)

	admin.Get("/tasks", h.getAllTasksAdmin)
	admin.Post("/tasks", h.createNewTask)

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask)
}

// @Summary      Get all tasks
// @Description  allows the admin to get all tasks.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param		 ?
// @Success      200  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /admin/tasks [get]
func (h *Handler) getAllTasksAdmin(ctx *fiber.Ctx) error {
	var filter domain.TaskFilter
	var tasks []*domain.Task

	err := ctx.BodyParser(&filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	tasks, err = h.AdminService.GetAllTasks(filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	type TaskRespose struct {
		Tasks []*domain.Task `json:"tasks"`
	}

	return ctx.Status(fiber.StatusOK).JSON(TaskRespose{
		Tasks: tasks,
	})
}

// @Summary      Update task
// @Description  allows the admin to update task.
// @Tags         UpdateTask
// @Accept       json
// @Produce      json
// @Param		 ?
// @Success      200  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /task/:id [patch]
func (h *Handler) updateTask(ctx *fiber.Ctx) error {

	taskId, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusConflict, err.Error())
	}

	var update domain.TaskUpdate

	err = ctx.BodyParser(&update)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	//TODO: VALIDATE taskData

	err = h.AdminService.UpdateTask(taskId, update)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
	})
}

// @Summary      Create new task
// @Description  allows the admin to create new tasks.
// @Tags         CreateTask
// @Accept       json
// @Produce      json
// @Param		 ?
// @Success      201  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/ [post]
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
