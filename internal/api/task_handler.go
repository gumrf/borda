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

func (h *Handler) getAllTasks(ctx *fiber.Ctx) error {

	var filter domain.TaskFilter
	var tasks []domain.TaskUserResponse

	err := ctx.BodyParser(&filter)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	id, _ := strconv.Atoi(ctx.Locals("userId").(string))

	tasks, err = h.UserService.ShowAllTasks(filter, id)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	type TaskRespose struct {
		Response []domain.TaskUserResponse
	}

	return ctx.Status(fiber.StatusOK).JSON(TaskRespose{
		Response: tasks,
	})
}

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
