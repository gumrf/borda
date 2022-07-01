package team

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"borda/internal/pkg/middleware"
	"borda/internal/pkg/response"
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
	data := request["payload"]

	switch method {
	case "create":
		_ = tc.teamService.CreateTeam(userId, CreatTeamRequest{Name: data})
	case "join":
		_ = tc.teamService.JoinTeam(userId, JoinTeamRequest{Token: data})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Status: strconv.Itoa(fiber.StatusBadRequest),
				Code:   response.IncorrectInputCode,
				Title:  "Неизвестный метод. Возможные варианты [create] и [join]'",
				Detail: err.Error(),
			},
		)
	}

	return c.SendStatus(fiber.StatusCreated)
}
