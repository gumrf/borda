package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTaskRoutes(router fiber.Router) {
	tasks := router.Group("/tasks", AuthRequired)
	tasks.Get("", h.getAllTasks)
	tasks.Post("", h.createNewTask)

	task := router.Group("/tasks/:id")
	task.Patch("", h.updateTask)

	task.Post("/submissions", h.createNewSubmission)
	task.Get("/submissions", h.getAllSubmissions)
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

	//TODO: VALIDATE task

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

	err = h.AdminUsecase.UpdateTask(taskId, taskData)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"error":   false,
		"message": "Successfully update task!",
	})
}

func (h *Handler) createNewSubmission(ctx *fiber.Ctx) error {
	var submission domain.SubmitTaskRequest
	err := ctx.BodyParser(&submission)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	//TODO: VALIDATE submission.Flag

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

func (h *Handler) getAllSubmissions(ctx *fiber.Ctx) error {

	var userTask domain.SubmitTaskRequest
	err := ctx.BodyParser(&userTask)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, err.Error())
	}

	var submissions []*domain.TaskSubmission
	submissions, err = h.UserUsecase.ShowAllSubmisiions(userTask)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusConflict, err.Error())
	}

	return ctx.JSON(fiber.Map{
		"error":       false,
		"submissions": submissions,
	})
}
