package api

import (
	"borda/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initScoreboardRoutes(router fiber.Router) {
	scoreboard := router.Group("/scoreboard")

	scoreboard.Get("", h.getScoreboard)
}

func (h *Handler) getScoreboard(c *fiber.Ctx) error {
	uc := usecase.NewUsecaseScoreboard(h.Repository.Teams, h.Repository.Tasks)

	response, err := uc.Execute()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
