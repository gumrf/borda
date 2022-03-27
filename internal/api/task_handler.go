package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", h.authRequired)
	tasks.Get("", h.getAllTasks)

	task := router.Group("/tasks/:id")

	task.Post("/submissions", h.createNewSubmission)
	task.Get("/submissions", h.getAllSubmissions)
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
	var tasks []domain.TaskUserResponse

	id, _ := strconv.Atoi(ctx.Locals("userId").(string))

	tasks, err := h.UserService.GetAllTasks(id)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	type TaskRespose struct {
		Tasks []domain.TaskUserResponse `json:"tasks"`
	}

	return ctx.Status(fiber.StatusOK).JSON(TaskRespose{
		Tasks: tasks,
	})
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
func (h *Handler) createNewSubmission(ctx *fiber.Ctx) error {
	var submission domain.SubmitTaskRequest
	err := ctx.BodyParser(&submission)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	if err := submission.Validate(); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input flag is invalid.")
	}

	var message string
	message, err = h.UserService.TryToSolveTask(submission)
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
func (h *Handler) getAllSubmissions(ctx *fiber.Ctx) error {

	var input domain.SubmitTaskRequest
	err := ctx.BodyParser(&input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	var submissions []*domain.TaskSubmission
	submissions, err = h.UserService.GetTaskSubmissions(input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusConflict, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"error":       false,
		"submissions": submissions,
	})
}
