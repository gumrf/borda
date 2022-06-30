package team

import (
	"borda/internal/pkg/middleware"
	"borda/internal/pkg/response"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TeamController struct {
	teamService *TeamService
}

func NewTeamController(ts *TeamService) *TeamController {
	return &TeamController{
		teamService: ts,
	}
}

func (tc *TeamController) InitTeamRoutes(router fiber.Router) {
	router.Post("/teams", middleware.AuthRequired, tc.joinTeam)
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
func (tc *TeamController) joinTeam(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int)

	var request map[string]string

	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Status: strconv.Itoa(fiber.StatusBadRequest),
				Code:   response.IncorrectInputCode,
				Title:  "Не могу распарсить жсон.",
				Detail: err.Error(),
			},
		)
	}

	method := request["method"]

	switch method {
	case "create":
		_ = tc.teamService.CreateTeam(userId, request["payload"])
	case "join":
		_ = tc.teamService.CreateTeam(userId, request["payload"])
	default:
		// TODO: send error
	}

	return c.SendStatus(fiber.StatusCreated)
}
