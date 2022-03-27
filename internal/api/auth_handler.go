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

func (h *Handler) handleSignUp(ctx *fiber.Ctx) error {
	var input domain.UserSignUpInput

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

func (h *Handler) handleSignIn(ctx *fiber.Ctx) error {
	var input domain.UserSignInInput

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
