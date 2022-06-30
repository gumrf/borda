package api

import (
	"borda/internal/domain"
	"borda/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAdminRoutes(router fiber.Router) {
	admin := router.Group("/admin", h.authRequired, h.adminPermissionRequired)

	tasks := admin.Group("/tasks")
	tasks.Get("", h.adminGetAllTasks)
	tasks.Post("", h.createNewTask)
	tasks.Patch("/:id", h.updateTask)
}

// @Summary      Get all tasks
// @Description  Get all tasks with admin access.
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {array}   domain.Task
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /admin/tasks [get]
func (h *Handler) adminGetAllTasks(c *fiber.Ctx) error {
	uc := usecase.NewAdminUsecaseGetTasks(h.Repository.Tasks)

	tasks, err := uc.Execute()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

// @Summary      Update task
// @Description  Update task.
// @Tags         Admin
// @Accept       json
// @Security     ApiKeyAuth
// @Produce      json
// @Param        task_id  path      int          true  "Task ID"
// @Param        task     body      domain.Task  true  "Task"
// @Success      200      string    OK
// @Failure      400,500  {object}  domain.ErrorResponse
// @Router       /admin/tasks/{task_id} [patch]
func (h *Handler) updateTask(c *fiber.Ctx) error {
	// Get task id from request url
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	// Get request payload
	var update domain.TaskUpdate
	if err := c.BodyParser(&update); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	if err := update.Validate("flag"); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, InvalidInputCode, "Input is invalid.", err.Error())
	}

	uc := usecase.NewAdminUsecaseUpdateTask(h.Repository.Tasks)

	updatedTask, err := uc.Execute(taskId, update)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

// @Summary      Create new task
// @Description  Create new task.
// @Tags         Admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        task     body      domain.Task  true  "Task"
// @Success      200      {object}  domain.Task
// @Failure      400,500  {object}  domain.ErrorResponse
// @Router       /admin/tasks [post]
func (h *Handler) createNewTask(c *fiber.Ctx) error {
	var task domain.Task

	if err := c.BodyParser(&task); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	//TODO: Fix validator: id, author id, points, isActive, isDisable
	//if err := task.Validate(); err != nil {
	//	return NewErrorResponse(ctx, fiber.StatusBadRequest,
	//		"Validation is not passed.", err.Error())
	//}

	uc := usecase.NewAdminUsecaseCreateTask(h.Repository.Tasks)

	createdTask, err := uc.Execute(task)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(createdTask)
}
