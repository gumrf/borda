package api

import (
	"borda/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users", h.authRequired, h.teamRequired)
	users.Get("/", h.getAllUsers)
	users.Get("/me", h.getCurentLoggedInUser)
	users.Get("/:id", h.getUser)
}

// @Summary      Get user profile
// @Description  Show curently logged in user profile.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  domain.PrivateUserProfileResponse
// @Failure      400      {object}  domain.ErrorResponse
// @Failure      404      {object}  domain.ErrorResponse
// @Failure      500      {object}  domain.ErrorResponse
// @Router       /users/me [get]
func (h *Handler) getCurentLoggedInUser(c *fiber.Ctx) error {
	id := c.Locals("userId").(int)

	uc := usecase.NewUserUsecaseGetProfile(h.Repository.Users, h.Repository.Teams)
	result, err := uc.Execute(id)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error occurred on the server", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary      Get user
// @Description  Show public user profile.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {array}   domain.PublicUserProfileResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /users/{user_id} [get]
func (h *Handler) getUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, IncorrectInputCode, "Input is incorrect", err.Error())
	}

	uc := usecase.NewUserUsecaseGetUsers(h.Repository.Users, h.Repository.Teams)

	result, err := uc.Execute(userId)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusConflict, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result[0])

}

// @Summary      Get all users
// @Description  Show all registered users.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200      {array}   domain.PublicUserProfileResponse
// @Failure      400,500  {object}  domain.ErrorResponse  "error"
// @Router       /users [get]
func (h *Handler) getAllUsers(c *fiber.Ctx) error {
	uc := usecase.NewUserUsecaseGetUsers(h.Repository.Users, h.Repository.Teams)

	result, err := uc.Execute()
	if err != nil {
		return NewErrorResponse(c, fiber.StatusInternalServerError, InternalServerErrorCode,
			"Internal error occurred on the server.", err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
