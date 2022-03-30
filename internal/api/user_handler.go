package api

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users", h.authRequired, h.checkUserInTeam)
	users.Get("", h.getAllUsers)
	users.Get("/:id", h.getUserById)
}

func (h *Handler) getUserById(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"error":  false,
		"userId": c.Params("id"),
	})
}

// @Summary      Get all users
// @Description  Show all registered users.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      201  {array}   domain.UserResponse
// @Failure      400  {object}  ErrorsResponse
// @Failure      404  {object}  ErrorsResponse
// @Failure      500  {object}  ErrorsResponse
// @Router       /users [get]
func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}
