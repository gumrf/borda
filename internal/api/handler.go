package api

import (
	_ "borda/docs"
	"borda/internal/config"
	"borda/internal/repository"
	"borda/internal/service"
	"fmt"
	"strconv"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	AuthService *service.AuthService
	Repository  *repository.Repository
}

func NewHandler(authService *service.AuthService, repository *repository.Repository) *Handler {
	return &Handler{
		AuthService: authService,
		Repository:  repository,
	}
}

func (h *Handler) Init(app *fiber.App) {
	app.Use(logger.New())

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			c.Path(): "pong",
			"time":   time.Now().Format(time.UnixDate),
		})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	h.initAuthRoutes(v1)
	h.initScoreboardRoutes(v1)

	// Everything defined bellow will require authorization
	v1.Use(jwtMiddleware.New(jwtMiddleware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return NewErrorResponse(c, fiber.StatusUnauthorized, "NOT_AUTHORIZED",
				"Authentication credentials are missing or invalid.",
				"Provide a properly configured and signed bearer token, and make sure that it has not expired.")
		},
		// TODO: DefineErrorHandler function
		SigningMethod: jwt.SigningMethodHS256.Name,
		SigningKey:    []byte(config.JWT().SigningKey),
		ContextKey:    "token",
	}))

	h.initUserRoutes(v1)
	h.initTaskRoutes(v1)
	h.initAdminRoutes(v1)
}

func (h *Handler) authRequired(c *fiber.Ctx) error {
	token := c.Locals("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Get user id, scope from claims
	id := claims["iss"].(string)
	scope := claims["scope"].([]interface{})

	// Convert user id from string to int
	intId, _ := strconv.Atoi(id)

	// Store user id, scope in context for the following routes
	c.Locals("userId", intId)
	c.Locals("scope", scope[0])

	fmt.Println("User ID: "+id+", Scope: ", scope[0])
	return c.Next()
}

func (h *Handler) checkUserInTeam(c *fiber.Ctx) error {
	id := c.Locals("userId").(int)

	teamId, ok := h.AuthService.VerifyUserTeam(id)
	if !ok {
		return NewErrorResponse(c, fiber.StatusForbidden,
			MissingTeamIdCode, "User is not a member of any group",
			"You tried to access a route that requires team id. Join a team before requesting this route.")
	}

	// Save team id to context
	c.Locals("teamId", teamId)

	return c.Next()
}

func (h *Handler) adminPermissionRequired(c *fiber.Ctx) error {
	scope := c.Locals("scope")
	if scope != "admin" {
		return NewErrorResponse(c, fiber.StatusForbidden, ForbiddenCode,
			"User does not have admin permission.",
			"Provide API key with right permission to access this route.",
		)
	}

	return c.Next()
}
