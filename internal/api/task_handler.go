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
	tasks.Post("/:id/flag", h.submitFlag)
	// tasks.Get("/:id/submissions", h.getAllSubmissions)
}

// @Summary      Get tasks
// @Description  Get all tasks available for user.
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200      {array}   domain.PublicTaskResponse
// @Failure      400,500  {object}  domain.ErrorResponse
// @Router       /tasks [get]
func (h *Handler) getAllTasks(c *fiber.Ctx) error {
	teamId := c.Locals("teamId").(int)

	uc := usecase.NewUserUsecaseGetTasks(h.Repository.Tasks, h.Repository.Teams)

	tasks, err := uc.Execute(teamId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}

// @Summary      Submit flag
// @Description  Try to solve task by sending flag.
// @Tags         Tasks
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        task_id  path      int                       true  "Task ID"
// @Param        flag     body      domain.SubmitFlagRequest  true  "Flag"
// @Success      201      string    OK
// @Failure      400,500  {object}  domain.ErrorResponse
// @Router       /tasks/{task_id}/flag [post]
func (h *Handler) submitFlag(c *fiber.Ctx) error {
	var input domain.SubmitFlagRequest
	if err := c.BodyParser(&input); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, InvalidInputCode, "Input is invalid.", err.Error())
	}

	taskId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())

	}

	uc := usecase.NewUserUsecaseSubmitFlag(h.Repository.Tasks)

	response, err := uc.Execute(domain.TaskSubmission{
		TaskId: taskId,
		TeamId: c.Locals("teamId").(int),
		UserId: c.Locals("userId").(int),
		Flag:   input.Flag,
	})
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
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
// func (h *Handler) getAllSubmissions(c *fiber.Ctx) error {
// 	// Get task from request url
// 	taskId, err := strconv.Atoi(c.Params("id"))
// 	if err != nil {
// 		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid.", err.Error())
// 	}

// 	// Get user id from context
// 	userId := c.Locals("userId").(int)
// 	uc := usecase.NewUserUsecaseGetAllSubmissions(h.Repository.Tasks)

// 	submissions, err := uc.Execute(taskId, userId)
// 	if err != nil {
// 		return NewErrorResponse(c, fiber.StatusConflict, "Error occurred on the server.", err.Error())
// 	}

// 	return c.Status(fiber.StatusOK).JSON(submissions)
// }
