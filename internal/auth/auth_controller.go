package auth

import (
	"borda/internal/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) InitRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/signup", ac.signUp)
	auth.Post("/signin", ac.signIn)
	auth.Post("/signout", ac.signOut)
	// TODO:
	// auth.Get("/verification", ac.verifyEmail)
}

// @Summary      Sign Up
// @Description  Register a new user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body     	RegistrationRequest   true  "Credentials"
// @Success      201          {string}  Created
// @Failure      400,500      {object}  response.ErrorResponse
// @Router       /auth/signup [post]
func (ac *AuthController) signUp(c *fiber.Ctx) error {
	var request RegistrationRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusBadRequest),
			Code:   response.IncorrectInputCode,
			Title:  "Input is incorrect",
			Detail: err.Error(),
		})
	}

	if err := request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusBadRequest),
			Code:   response.InvalidInputCode,
			Title:  "Incorrect request body",
			Detail: err.Error(),
		})
	}

	if err := ac.authService.Register(request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusInternalServerError),
			Code:   response.InternalServerErrorCode,
			Title:  "Internal error occurred on the server.",
			Detail: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary      Sign In
// @Description  Authorize user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      AuthorizationRequest  true  "Credentials"
// @Success      200          {object}
// @Failure      default      {object}  response.ErrorResponse
// @Router       /auth/sign-in [post]
func (ac *AuthController) signIn(c *fiber.Ctx) error {
	var request AuthorizationRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusBadRequest),
			Code:   response.IncorrectInputCode,
			Title:  "Input is incorrect",
			Detail: err.Error(),
		})
	}

	if err := request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusBadRequest),
			Code:   response.InvalidInputCode,
			Title:  "Incorrect request body",
			Detail: err.Error(),
		})
	}

	token, err := ac.authService.Authorize(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusInternalServerError),
			Code:   response.InternalServerErrorCode,
			Title:  "Internal error occurred on the server.",
			Detail: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthController) signOut(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(response.ErrorResponse{
		Status: strconv.Itoa(fiber.StatusNotImplemented),
		Code:   "NOT_IMPLEMENTED",
		Title:  "Api route is not implemented.",
	})
}
