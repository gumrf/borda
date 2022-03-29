package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", h.authRequired, h.checkUserInTeam)
	tasks.Get("", h.getAllTasks)

	tasks.Post("/:id/submissions", h.submitFlag)
	tasks.Get("/:id/submissions", h.getAllSubmissions)
}

// @Summary      Get all tasks
// @Description  allows the user to get tasks by filter.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Success      200  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/ [get]
func (h *Handler) getAllTasks(ctx *fiber.Ctx) error {
	id := ctx.Locals("userId").(int)

	tasks, err := h.UserService.GetAllTasks(id)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Error occurred on the server.", err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"tasks": tasks})
}

// @Summary      Create new submission
// @Description  allows the user to create submission.
// @Tags         Submissions
// @Accept       json
// @Produce      json
// @Param		 id   path      int  true  "Task ID"
// @Success      201  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/:id/submissions [post]
func (h *Handler) submitFlag(c *fiber.Ctx) error {
	var submission domain.SubmitFlagRequest
	if err := c.BodyParser(&submission); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := submission.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.", err.Error())
	}

	// Get task from request url
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.", err.Error())
	}

	if err := h.UserService.SolveTask(domain.TaskSubmission{
		TaskId: taskId,
		TeamId: c.Locals("teamId").(int),
		UserId: c.Locals("userId").(int),
		Flag:   submission.Flag,
	}); err != nil {
		return NewErrorResponse(c,
			fiber.StatusConflict, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// @Summary      Get all submission
// @Description  allows the user to get all submissions for task.
// @Tags         Submissions
// @Accept       json
// @Produce      json
// @Param		 id   path      int  true  "Task ID"
// @Success      201  {object}  TaskResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tasks/:id/submissions [get]
func (h *Handler) getAllSubmissions(c *fiber.Ctx) error {
	// Get task from request url
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.", err.Error())
	}

	// Get user id from context
	userId := c.Locals("userId").(int)

	submissions, err := h.UserService.GetTaskSubmissions(taskId, userId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusConflict, "", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"submissions": submissions})
}
