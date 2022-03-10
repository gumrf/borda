package api

// import (
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// )

// func RegisterRoutes(app *fiber.App) {
// 	app.Use(logger.New())
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"message": "OK",
// 			"time":    time.Now().Format(time.UnixDate),
// 		})
// 	})

// 	api := app.Group("/api")
// 	v1 := api.Group("/v1")

// 	initHealthCheckRoute(v1)
// 	initUserRoutes(v1)
// 	initTaskRoutes(v1)
// }
