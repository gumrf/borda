package api

import (
	"borda/internal/domain"
	"borda/internal/usecase"
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
// @Description  Get tasks.
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  domain.UserTaskResponse
// @Failure      400      {object}  ErrorsResponse
// @Failure      404      {object}  ErrorsResponse
// @Failure      500      {object}  ErrorsResponse
// @Router       /tasks [get]
func (h *Handler) getAllTasks(ctx *fiber.Ctx) error {
	id := ctx.Locals("teamId").(int)

	uc := usecase.NewUsecaseGetAllTasks(h.Repository.Tasks, h.Repository.Teams)

	tasks, err := uc.Execute(id)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Error occurred on the server.", err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"tasks": tasks})
}

// @Summary      Submit flag
// @Description  Create new flag submission.
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        flag     body      domain.SubmitFlagRequest  true  "Flag"
// @Param        task_id  path      int                       true  "Task ID"
// @Success      201      string    OK
// @Failure      400      {object}  ErrorsResponse
// @Failure      404      {object}  ErrorsResponse
// @Failure      500      {object}  ErrorsResponse
// @Router       /tasks/{task_id}/submissions [post]
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

	uc := usecase.NewUsecaseSubmitFlag(h.Repository.Tasks)

	if err := uc.Execute(domain.TaskSubmission{
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
// @Description  Get all submissions for task.
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Produce      json
// @Param        task_id  path      int  true  "Task ID"
// @Success      201      {object}  domain.TaskSubmission
// @Failure      400  {object}  ErrorsResponse
// @Failure      404  {object}  ErrorsResponse
// @Failure      500  {object}  ErrorsResponse
// @Router       /tasks/{task_id}/submissions [get]
func (h *Handler) getAllSubmissions(c *fiber.Ctx) error {
	// Get task from request url
	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.", err.Error())
	}

	// Get user id from context
	userId := c.Locals("userId").(int)
	uc := usecase.NewUsecaseGetAllSubmissions(h.Repository.Tasks)

	submissions, err := uc.Execute(taskId, userId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusConflict, "", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"submissions": submissions})
}
