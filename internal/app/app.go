package app

import (
	_ "borda/api"
	"borda/internal/config"
	"borda/internal/pkg/response"
	"borda/pkg/log"
	"borda/pkg/pg"

	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	loggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	jwtMiddleware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// type App struct {
// 	mutex  sync.Mutex
// 	server *fiber.App
// 	config Config
// }

// @title                       Borda API
// @version                     0.1
// @description                 REST API for Borda.
// @host                        localhost:8080
// @BasePath                    /api/v1
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// Run initialize the application
func Run() {
	// ctx := context.Background()

	config := config.NewConfig()
	logger := log.New()

	logger.Info(config.String())

	logger.Infof("Connecting to Postgres: %s", config.DBURL())
	db, err := pg.Connect(config.DBURL())
	if err != nil {
		logger.Fatalf("Failed to connect to Postgres: %v", err)
	}
	logger.Info("Connected to Postgres: ", config.DBURL())

	migrateError := pg.Migrate(config.DBURL(), config.MigrationsDir(), 2)
	if migrateError != nil {
		logger.Fatalf("Failed to run migrations: %v", migrateError)
	}

	// repository := repository.NewRepository(db)

	// authService := service.NewAuthService(repository.Users, repository.Teams,
	// 	hash.NewSHA1Hasher(config.PasswordSalt()),
	// )

	// userRepository := user.NewUserRepository(db)
	// teamRepository := team.NewTeamRepository(db)

	// userService := user.NewUserService(userRepository, teamRepository)
	// userController := user.NewUserController(userService)

	// authService := auth.NewAuthService(userRepository, teamRepository, hash.NewSHA1Hasher(config.Salt()), config.JWT())
	// authController := auth.NewAuthController(authService)

	app := fiber.New()
	app.Use(loggerMiddleware.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api/v1")

	api.Get("/status", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK",
			"time":   time.Now().Format(time.UnixDate),
		})
	})

	// authController.InitRoutes(api)

	// Everything defined bellow will require authorization
	api.Use(jwtMiddleware.New(jwtMiddleware.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Status: strconv.Itoa(fiber.StatusUnauthorized),
				Code:   response.NotAuthorizedCode,
				Title:  "Authentication credentials are missing or invalid.",
				Detail: "Provide a properly configured and signed bearer token, and make sure that it has not expired.",
			})
		},
		// TODO: DefineErrorHandler function
		SigningMethod: jwt.SigningMethodHS256.Name,
		SigningKey:    []byte(config.JWT().SigningKey),
		ContextKey:    "token",
	}))

	api.Static("/files", "./_content", fiber.Static{Download: true})

	// userController.InitRoutes(api)
	// teamController.InitRoutes(api)

	app.Use(
		func(c *fiber.Ctx) error {
			// Return HTTP 404 status and JSON response.
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
				Status: strconv.Itoa(fiber.StatusNotFound),
				Code:   "NOT_FOUND",
				Title:  fmt.Sprintf("Cannot %s %s", c.Method(), c.BaseURL()+c.OriginalURL()),
			})
		},
	)

	// Catch OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit

		logger.Info("Received an interrupt signal, shutdown")

		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Errorf("Server shutdown failed: %v", err)
		}

		db.Close()
	}()

	if err := app.Listen(":" + config.GetString("app.port")); err != nil {
		logger.Errorf("Failed to run BordaApplication: %v", err)
	}
}
