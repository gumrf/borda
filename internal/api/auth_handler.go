package api

import (
	"borda/internal/domain"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/signUp", h.handleSignUp)
	auth.Post("/signIn", h.handleSignIn)
	auth.Post("/signOut", h.handleSignOut)
}

func (h *Handler) handleSignUp(ctx *fiber.Ctx) error {
	var input domain.UserSignUpInput

	err := ctx.BodyParser(&input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.")
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.")
	}

	err = h.AuthService.SignUp(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignUp
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, fmt.Sprintf("Error occurred on the server. Error: %s", err.Error()),
		)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *Handler) handleSignIn(ctx *fiber.Ctx) error {
	var input domain.UserSignInInput

	err := ctx.BodyParser(&input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.")
	}

	if err := input.Validate(); err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.")
	}

	token, err := h.AuthService.SignIn(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignIn
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, err.Error(),
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) handleSignOut(ctx *fiber.Ctx) error {
	return NewErrorResponse(ctx,
		fiber.StatusNotImplemented, "Api route is not implemented.")
}
