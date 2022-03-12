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

func (h *Handler) handleSignUp(c *fiber.Ctx) error {
	var input domain.UserSignUpInput

	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "400",
				Code:   "BAD_REQUEST",
			}),
		)
	}

	// err = input.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "400",
				Code:   "BAD_REQUEST",
				// TODO: title, details content depends on return error
				Title: "Input data is invalid.",
			}),
		)
	}

	err = h.AuthService.SignUp(input)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignUp
		return c.Status(500).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "500",
				Code:   "INTERNAL_SERVER_ERROR",
			}),
		)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) handleSignIn(c *fiber.Ctx) error {
	var user domain.UserSignInInput

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "400",
				Code:   "BAD_REQUEST",
			}),
		)
	}

	// VALIDATE USERNAME AND PASSWORD
	//if err := Validate(username, pssword); err != nil {
	//	return c.Status(400).JSON(
	//		NewAPIErrorResponse(ErrorObject{
	//			Status: "400",
	//			Code:   "BAD_PSWD/UNAME",
	//		}),
	//	)
	//}

	token, err := h.AuthService.SignIn(user.Username, user.Password)
	if err != nil {
		// TODO: send specific error depending on error type returned by SignIn
		return c.Status(500).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "500",
				Code:   "INTERNAL_SERVER_ERROR",
			}),
		)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) handleSignOut(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(
		NewAPIErrorResponse(ErrorObject{
			Status: "501",
			Code:   "NOT_IMPLEMENTED",
			Title:  "Api route is not implemented.",
		}),
	)
}
