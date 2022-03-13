package api

import (
	"borda/internal/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	AuthService *services.AuthService
}

func NewHandler(authService *services.AuthService) *Handler {
	return &Handler{
		AuthService: authService,
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
	h.initUserRoutes(v1)
	h.initTaskRoutes(v1)
}

func AuthRequired(c *fiber.Ctx) error {

	return c.Next()
}
