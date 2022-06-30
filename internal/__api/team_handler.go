package api

import (
	"borda/internal/domain"
	"borda/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTeamRoutes(router fiber.Router) {
	team := router.Group("/teams", h.authRequired)
	team.Post("", h.joinTeam)
}

// @Summary      Join or create team
// @Description  Join team by token or create a new team.
// @Tags         Team
// @Accept       json
// @Produce      json
// @Param        TeamInput  body      domain.TeamInput  true  "Team Input"
// @Success      201          {string}  Created
// @Failure      400,500      {object}  domain.ErrorResponse
// @Router       /team [post]
func (h *Handler) joinTeam(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	var input domain.AttachTeamInput

	if err := c.BodyParser(&input); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, InvalidInputCode, "Input is invalid.", err.Error())
	}

	uc := usecase.NewUserUsecaseJoinTeam(h.Repository.Users, h.Repository.Teams)

	if err := uc.Execute(userId, input.Method, input.Attribute); err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Interal error occurred on the server", err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}
