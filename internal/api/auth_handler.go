package api

import (
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
	return c.JSON(fiber.Map{
		"token": "JWT-Token",
	})
}

func (h *Handler) handleSignIn(c *fiber.Ctx) error {
	// Get accountName, password -> bad request
	// Validate SignInInput -> bad request
	// Generate new token -> internal errror
	// return new token
	return c.JSON(fiber.Map{
		"error": false,
		"token": "JWT-Token",
	})
}

func (h *Handler) handleSignOut(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error": false,
	})
}
