package api

import (
	"borda/internal/domain"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Response(c *fiber.Ctx, status int, body interface{}) error {
	return c.Status(status).JSON(body)
}

func NewErrorResponse(c *fiber.Ctx, status int, code string, title string, detail ...string) error {
	return c.Status(status).JSON(
		domain.ErrorResponse{
			Status: strconv.Itoa(status),
			Code:   code,
			Title:  title,
			Detail: func(detail []string) string {
				if len(detail) > 0 {
					return detail[0]
				}
				return ""
			}(detail),
		},
	)
}

var (
	MissingTeamIdCode       = "MISSING_TEAM_ID"
	ForbiddenCode           = "FORBIDDEN"
	IncorrectInputCode      = "INCORRECT_INPUT"
	InvalidInputCode        = "INVALID_INPUT"
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
)
