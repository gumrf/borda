package api

import (
	"borda/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.handleSignUp)
	auth.Post("/sign-in", h.handleSignIn)
	auth.Post("/sign-out", h.handleSignOut)
}

// @Summary      Sign Up
// @Description  Create a new user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.SignUpInput  true  "Credentials"
// @Success      201          {string}  Created
// @Failure      400,500      {object}  domain.ErrorResponse
// @Router       /auth/sign-up [post]
func (h *Handler) handleSignUp(c *fiber.Ctx) error {
	var input domain.SignUpInput

	if err := c.BodyParser(&input); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, InvalidInputCode, "Input is invalid.", err.Error())
	}

	if err := h.AuthService.SignUp(input); err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary      Sign In
// @Description  Sign in into account.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.SignInInput  true  "Credentials"
// @success      200          {object}  domain.SignInResponse
// @Failure      400,500      {object}  domain.ErrorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) handleSignIn(c *fiber.Ctx) error {
	var input domain.SignInInput

	if err := c.BodyParser(&input); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, InvalidInputCode, "Input is invalid.", err.Error())
	}

	token, err := h.AuthService.SignIn(input)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) handleSignOut(c *fiber.Ctx) error {
	return NewErrorResponse(c, fiber.StatusNotImplemented, "NOT_IMPLEMENTED", "Api route is not implemented.")
}
