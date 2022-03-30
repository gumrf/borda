package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users", h.authRequired, h.checkUserInTeam)
	users.Get("", h.getAllUsers)

	user := users.Group("/:id")
	user.Get("", h.getUserById)
}

func (h *Handler) getUserById(c *fiber.Ctx) error {

	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Input is invalid ", err.Error())
	}

	user, err := h.UserService.GetUserById(userId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server. ", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user_profile": user})
}

func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": users})
}
