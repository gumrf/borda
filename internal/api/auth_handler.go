package api

import (
	"borda/internal/domain"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initAuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	// registrinton
	auth.Post("/signUp", h.handleSignUp)
	auth.Post("/signIn", h.handleSignIn)
	auth.Post("/signOut", h.handleSignOut)
}

func (h *Handler) handleSignUp(c *fiber.Ctx) error {
	var user domain.User

	if err := c.BodyParser(user); err != nil {
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

	id, err := h.AuthService.SignUp(user)
	if err != nil {
		return c.Status(500).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "500",
				Code:   "INTERNAL_SERVER_ERROR",
				Detail: "BAD PSWD OR USERNAME",
			}),
		)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"id":    id,
	})
}

func (h *Handler) handleSignIn(c *fiber.Ctx) error {
	var user domain.User

	if err := c.BodyParser(user); err != nil {
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
		return c.Status(500).JSON(
			NewAPIErrorResponse(ErrorObject{
				Status: "500",
				Code:   "INTERNAL_SERVER_ERROR",
			}),
		)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"token": token,
	})
}

func (h *Handler) handleSignOut(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
	})
}
