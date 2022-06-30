package user

import (
	"borda/internal/pkg/middleware"
	"borda/internal/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService *UserService
}

func NewUserController(userService *UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (uc *UserController) InitRoutes(router fiber.Router) {
	users := router.Group("/users", middleware.AuthRequired)
	users.Get("/", uc.getAllUsers)
	users.Get("/:userId", uc.getUser)
	users.Get("/me", uc.getUserProfile)
}

type ErrorResponse response.ErrorResponse

// @Summary      Get user profile
// @Description  Show curently logged in user profile.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200       {object}  UserProfileResponse
// @Failure      default   {object}  ErrorResponse
// @Router       /users/me [get]
func (uc *UserController) getUserProfile(c *fiber.Ctx) error {
	userId := c.Locals("USER_ID").(int)

	response, err := uc.userService.GetUserProfile(userId)
	if err != nil {
		return c.Status(fiber.StatusNotImplemented).JSON(
			ErrorResponse{
				Status: strconv.Itoa(fiber.StatusNotImplemented),
				Code:   "NOT_IMPLEMENTED",
				Title:  "Api route is not implemented.",
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// @Summary      Get user
// @Description  Get user by id.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Param        userId  path      int  true  "User ID"
// @Success      200      {object}  UserResponse
// @Failure      default  {object}  ErrorResponse
// @Router       /users/{userId} [get]
func (uc *UserController) getUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorResponse{
				Status: strconv.Itoa(fiber.StatusBadRequest),
				Code:   response.IncorrectInputCode,
				Title:  "Input is incorrect",
				Detail: "Missing or wrong url param id. Example: /users/{123",
			},
		)
	}

	userResponse, err := uc.userService.GetUser(userId)
	if err != nil {
		// What kind of error???
		return c.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status: strconv.Itoa(fiber.StatusInternalServerError),
				Code:   response.InternalServerErrorCode,
				Title:  "Что-то пошло не так(((",
				Detail: err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
}

// @Summary      Get all users
// @Description  Show all registered users.
// @Tags         Users
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200      {array}   UserResponse
// @Failure      default  {object}  ErrorResponse
// @Router       /users [get]
func (uc *UserController) getAllUsers(c *fiber.Ctx) error {
	usersResponse, err := uc.userService.GetAllUsers()
	if err != nil {
		// What kind of error???
		return c.Status(fiber.StatusInternalServerError).JSON(
			ErrorResponse{
				Status: strconv.Itoa(fiber.StatusInternalServerError),
				Code:   response.InternalServerErrorCode,
				Title:  "Что-то пошло не так(((",
				Detail: err.Error(),
			},
		)
	}

	return c.Status(fiber.StatusOK).JSON(usersResponse)
}
