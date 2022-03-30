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
// @Failure      400          {object}  ErrorsResponse
// @Failure      404          {object}  ErrorsResponse
// @Failure      500          {object}  ErrorsResponse
// @Router       /auth/sign-up [post]
func (h *Handler) handleSignUp(ctx *fiber.Ctx) error {
	var input domain.SignUpInput

	if err := ctx.BodyParser(&input); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.", " Validation is not passed: "+err.Error())
	}

	if err := h.AuthService.SignUp(input); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, "Error occurred on the server.", err.Error(),
		)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

// @Summary      Sign In
// @Description  User sign in.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.SignInInput  true  "Credentials"
// @success      200          {string}  token
// @Failure      400          {object}  ErrorsResponse
// @Failure      404          {object}  ErrorsResponse
// @Failure      500          {object}  ErrorsResponse
// @Router       /auth/sign-in [post]
func (h *Handler) handleSignIn(ctx *fiber.Ctx) error {
	var input domain.SignInInput

	if err := ctx.BodyParser(&input); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.", err.Error())
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.", err.Error())
	}

	token, err := h.AuthService.SignIn(input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, "Error occurred on the server.", err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) handleSignOut(ctx *fiber.Ctx) error {
	return NewErrorResponse(ctx,
		fiber.StatusNotImplemented, "Api route is not implemented.")
}
