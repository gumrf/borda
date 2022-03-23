package api

import (
	"borda/internal/domain"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.handleSignUp)
	auth.Post("/sign-in", h.handleSignIn)
	auth.Post("/sign-out", h.handleSignOut)
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

	err = h.AuthService.DataVerification(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignUp
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, fmt.Sprintf("Error occurred on the server. Error: %s", err.Error()),
		)
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

	// TODO: Input validation

	token, err := h.AuthService.SignIn(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignIn
		return NewErrorResponse(ctx,
			fiber.StatusInternalServerError, err.Error(),
		)
	}

	//Здесь следует проверка есть ли у чела команда, если да то проходи и на токен если нет то создай команду сучка

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) handleSignOut(ctx *fiber.Ctx) error {
	return NewErrorResponse(ctx,
		fiber.StatusNotImplemented, "Api route is not implemented.")
}
