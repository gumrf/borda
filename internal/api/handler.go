package api

import (
	"borda/internal/config"
	"borda/internal/service"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type Handler struct {
	AuthService  *service.AuthService
	UserService  *service.UserService
	AdminService *service.AdminService
}

func NewHandler(authService *service.AuthService, userService *service.UserService,
	adminService *service.AdminService) *Handler {
	return &Handler{
		AuthService:  authService,
		UserService:  userService,
		AdminService: adminService,
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

	api := app.Group("/api")
	v1 := api.Group("/v1")

	h.initAuthRoutes(v1)

	// Everything defined bellow will require authorization
	v1.Use(jwtMiddleware.New(jwtMiddleware.Config{
		SigningMethod: jwt.SigningMethodHS256.Name,
		SigningKey:    []byte(config.JWT().SigningKey),
		ContextKey:    "token",
	}))

	h.initUserRoutes(v1)
	h.initTaskRoutes(v1)
}

func (h *Handler) authRequired(c *fiber.Ctx) error {
	token := c.Locals("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Get user id, scope from claims
	id := claims["iss"].(string)
	scope := claims["scope"].([]interface{})

	// Store user id, scope in context for the following routes
	c.Locals("userId", id)
	c.Locals("scope", scope[0])

	fmt.Println("User ID: "+id+", Scope: ", scope[0])
	return c.Next()
}

func (h *Handler) CheckUserInTeam(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Locals("userId").(string))

	if !h.UserService.IsUserInTeam(id) {
		return NewErrorResponse(c, fiber.StatusForbidden, "Authorization is not completed.") // Add details for err
	}

	return c.Next()
}

func (h *Handler) adminPermissionRequired(c *fiber.Ctx) error {
	scope := c.Locals("scope")
	if scope != "admin" {
		return NewErrorResponse(c, fiber.StatusForbidden, "You are not allowed to access resource. Ask for admin permission")
	}

	return c.Next()
}
