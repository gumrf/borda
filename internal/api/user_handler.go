package api

import (
	"borda/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users", h.authRequired, h.checkUserInTeam)
	users.Get("/", h.getAllUsers)
	users.Get("/me", h.getCurentLogetInUser)
	users.Get("/:id", h.getUser)
}

// @Summary      Get curent loget in user
// @Description  Show page of curent loget in user.
// @Tags         User
// @Security     ApiKeyAuth
// @Produce      json
// @Success      201  {array}   domain.UserProfileResponse
// @Failure      400  {object}  ErrorsResponse
// @Failure      404  {object}  ErrorsResponse
// @Failure      500  {object}  ErrorsResponse
// @Router       /users/me [get]
func (h *Handler) getCurentLogetInUser(c *fiber.Ctx) error {
	id := c.Locals("userId").(int)

	uc := usecase.NewUsecaseGetUser(h.Repository.Users, h.Repository.Teams)
	result, err := uc.Execute(id, true)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"profile": result,
	})
}

// @Summary      Get user
// @Description  Show page of user.
// @Tags         User
// @Security     ApiKeyAuth
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Success      201  {array}   domain.UserProfileResponse
// @Failure      400  {object}  ErrorsResponse
// @Failure      404  {object}  ErrorsResponse
// @Failure      500  {object}  ErrorsResponse
// @Router       /users/{user_id} [get]
func (h *Handler) getUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on convertation id in int", err.Error())
	}

	uc := usecase.NewUsecaseGetUser(h.Repository.Users, h.Repository.Teams)

	result, err := uc.Execute(userId, false)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"profile": result,
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

	uc := usecase.NewUsecaseGetAllUsers(h.Repository.Users, h.Repository.Teams)

	result, err := uc.Execute()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"users": result})
}
