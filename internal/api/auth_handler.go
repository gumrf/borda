package api

import (
	"borda/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/signUp", h.handleSignUp)
	auth.Post("/signIn", h.handleSignIn)
	auth.Post("/signOut", h.handleSignOut)
}

// api/v1/auth/signUp?team[create]=teamName
// api/v1/auth/signUp?team[join]=token

func (h *Handler) handleSignUp(ctx *fiber.Ctx) error {
	var input domain.UserSignUpInput

	err := ctx.BodyParser(&input)
	if err != nil {
		return NewErrorResponse(ctx,
			fiber.StatusBadRequest, "Input data is invalid.")
	}

	// TODO: Input validation

	err = h.AuthService.SignUp(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignUp
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, "Error occurred on the server.",
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

	// TODO: Input validation

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
